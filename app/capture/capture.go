package capture

import (
	"auto-record/app/window"
	"context"
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
	"github.com/pkg/errors"
	"github.com/tsqqqqqq/GoRecord/winapi"
	"github.com/tsqqqqqq/GoRecord/winapi/dx11"
	"github.com/tsqqqqqq/GoRecord/winapi/winrt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"os/signal"
	cr "runtime"
	"syscall"
	"time"
	"unsafe"
)

type Capture struct {
	ctx context.Context
}

// CaptureCount 每当有一个键盘鼠标得输入事件的时候， 按照时间戳截图一张当前窗口句柄。。。
var CaptureCount int = 0
var ImageChan chan string

const ext = "png"

func init() {
	if ImageChan == nil {
		ImageChan = make(chan string, 2)
	}
}

func NewCapture(ctx context.Context) *Capture {
	return &Capture{ctx}
}

type CaptureHandler struct {
	device                 *winrt.IDirect3DDevice
	deviceDx               *dx11.ID3D11Device
	graphicsCaptureItem    *winrt.IGraphicsCaptureItem
	framePool              *winrt.IDirect3D11CaptureFramePool
	graphicsCaptureSession *winrt.IGraphicsCaptureSession
	framePoolToken         *winrt.EventRegistrationToken
	isRunning              bool
}

func (c *CaptureHandler) StartCapture(hwnd win.HWND) error {
	type resultAttr struct {
		err error
	}

	var result = make(chan resultAttr)

	go func() {
		cr.LockOSThread()
		defer cr.UnlockOSThread()

		// Initialize Windows Runtime
		err := winrt.RoInitialize(winrt.RO_INIT_MULTITHREADED)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "RoInitialize")}
			return
		}
		defer winrt.RoUninitialize()

		// Create capture device
		var featureLevels = []dx11.D3D_FEATURE_LEVEL{
			dx11.D3D_FEATURE_LEVEL_11_0,
			dx11.D3D_FEATURE_LEVEL_10_1,
			dx11.D3D_FEATURE_LEVEL_10_0,
			dx11.D3D_FEATURE_LEVEL_9_3,
			dx11.D3D_FEATURE_LEVEL_9_2,
			dx11.D3D_FEATURE_LEVEL_9_1,
		}

		err = dx11.D3D11CreateDevice(
			nil, dx11.D3D_DRIVER_TYPE_HARDWARE, 0, dx11.D3D11_CREATE_DEVICE_BGRA_SUPPORT|dx11.D3D11_CREATE_DEVICE_DEBUG,
			&featureLevels[0], len(featureLevels),
			dx11.D3D11_SDK_VERSION, &c.deviceDx, nil, nil,
		)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "D3DCreateDevice")}
			return
		}
		defer c.deviceDx.Release()

		// Query interface of DXGIDevice
		var dxgiDevice *dx11.IDXGIDevice
		err = c.deviceDx.PutQueryInterface(dx11.IDXGIDeviceID, &dxgiDevice)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "PutQueryInterface")}
			return
		}

		var deviceRT *ole.IInspectable

		// convert D3D11Device(Dx11) to Direct3DDevice(WinRT)
		err = dx11.CreateDirect3D11DeviceFromDXGIDevice(dxgiDevice, &deviceRT)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "CreateDirect3D11DeviceFromDXGIDevice")}
			return
		}
		defer deviceRT.Release()

		// Query interface of IDirect3DDevice
		err = deviceRT.PutQueryInterface(winrt.IDirect3DDeviceID, &c.device)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "QueryInterface: IDirect3DDeviceID")}
			return
		}
		defer c.device.Release()

		// Create Capture Settings
		factory, err := ole.RoGetActivationFactory(winrt.GraphicsCaptureItemClass, winrt.IGraphicsCaptureItemInteropID)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "RoGetActivationFactory: IGraphicsCaptureItemID")}
			return
		}
		defer factory.Release()

		var interop *winrt.IGraphicsCaptureItemInterop
		err = factory.PutQueryInterface(winrt.IGraphicsCaptureItemInteropID, &interop)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "QueryInterface: IGraphicsCaptureItemInteropID")}
			return
		}
		defer interop.Release()

		var captureItemDispatch *ole.IInspectable

		// Capture for the window specified
		err = interop.CreateForWindow(hwnd, winrt.IGraphicsCaptureItemID, &captureItemDispatch)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "CreateForMonitor")}
			return
		}
		defer captureItemDispatch.Release()

		// Capture for the monitor specified
		//var hmoni = win.MonitorFromWindow(hwnd, win.MONITORINFOF_PRIMARY)
		//
		//err = interop.CreateForMonitor(hmoni, winrt.IGraphicsCaptureItemID, &captureItemDispatch)
		//if err != nil {
		//	result <- resultAttr{errors.Wrap(err, "CreateForMonitor")}
		//	return
		//}
		//defer captureItemDispatch.Release()

		// Get Interface of IGraphicsCaptureItem
		// 解释一下下面的代码
		// 1. captureItemDispatch是一个IInspectable对象
		err = captureItemDispatch.PutQueryInterface(winrt.IGraphicsCaptureItemID, &c.graphicsCaptureItem)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "PutQueryInterface captureItemDispatch")}
			return
		}

		// Get Capture objects size
		size, err := c.graphicsCaptureItem.Size()
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "Size")}
			return
		}

		// Get object of Direct3D11CaptureFramePoolClass
		ins, err := ole.RoGetActivationFactory(winrt.Direct3D11CaptureFramePoolClass, winrt.IDirect3D11CaptureFramePoolStaticsID)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "RoGetActivationFactory: IDirect3D11CaptureFramePoolStatics Class Instance")}
			return
		}
		defer ins.Release()

		// Get Interface of Direct3D11CaptureFramePoolClass
		var framePoolStatic *winrt.IDirect3D11CaptureFramePoolStatics2
		err = ins.PutQueryInterface(winrt.IDirect3D11CaptureFramePoolStatics2ID, &framePoolStatic)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "PutQueryInterface: IDirect3D11CaptureFramePoolStaticsID")}
			return
		}
		defer framePoolStatic.Release()

		// Create frame pool
		c.framePool, err = framePoolStatic.CreateFreeThreaded(c.device, winrt.DirectXPixelFormat_B8G8R8A8UIntNormalized, 1, size)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "CreateFramePool")}
			return
		}

		// Set frame settings
		var eventObject = NewDirect3D11CaptureFramePool(c.onFrameArrived)
		c.framePoolToken, err = c.framePool.AddFrameArrived(unsafe.Pointer(eventObject))
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "AddFrameArrived")}
			return
		}
		defer eventObject.Release()

		c.graphicsCaptureSession, err = c.framePool.CreateCaptureSession(c.graphicsCaptureItem)
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "CreateCaptureSession")}
			return
		}
		//fmt.Println(c.graphicsCaptureSession.IsSupport())
		defer c.graphicsCaptureSession.Release()

		// Start capturing
		err = c.graphicsCaptureSession.StartCapture()
		if err != nil {
			result <- resultAttr{errors.Wrap(err, "StartCapture")}
			return
		}

		c.isRunning = true

		result <- resultAttr{nil}

		for c.isRunning {
			time.Sleep(time.Second)
		}
	}()

	var res = <-result
	close(result)

	return res.err
}

func (c *CaptureHandler) onFrameArrived(this_ *uintptr, sender *winrt.IDirect3D11CaptureFramePool, args *ole.IInspectable) uintptr {
	_ = (*Direct3D11CaptureFramePool)(unsafe.Pointer(this_))
	frame, err := sender.TryGetNextFrame()
	frame.Process()
	if err != nil {
		os.Stderr.Write([]byte("Error: TryGetNextFrame: " + err.Error()))
		return 0
	}
	return 0
}

func (c *CaptureHandler) Close() error {
	if !c.isRunning {
		return nil
	}

	if c.framePool != nil {
		err := c.framePool.RemoveFrameArrived(c.framePoolToken)
		if err != nil {
			return errors.Wrap(err, "RemoveFrameArrived")
		}

		var closable *winrt.IClosable
		err = c.framePool.PutQueryInterface(winrt.IClosableID, &closable)
		if err != nil {
			return errors.Wrap(err, "PutQueryInterface: graphicsCaptureSession")
		}
		defer closable.Release()

		closable.Close()

		c.framePool = nil
	}

	var closable *winrt.IClosable
	err := c.graphicsCaptureSession.PutQueryInterface(winrt.IClosableID, &closable)
	if err != nil {
		return errors.Wrap(err, "PutQueryInterface: graphicsCaptureSession")
	}
	defer closable.Release()

	closable.Close()

	c.graphicsCaptureItem = nil
	c.isRunning = false

	return nil
}

func WindowCapture(input chan string) {
	var rdHwnd win.HWND
	for _ = range input {
		// TODO 这里还需要优化以下 最后在根据操作生成视频
		currentWindow := window.CurrentWindow
		robotgo.SetActive(robotgo.GetHandPid(currentWindow.Pid))
		rdHwnd = winapi.FindWindow(nil, winapi.MustUTF16PtrFromString(currentWindow.Title))
		if rdHwnd == 0 {
			win.MessageBox(0, winapi.MustUTF16PtrFromString("Could not find window"), winapi.MustUTF16PtrFromString("RDP Relative Input"), win.MB_ICONERROR)
		}
		handler := new(CaptureHandler)
		err := handler.StartCapture(rdHwnd)
		defer handler.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sc
	}

	//for _ = range input {
	//	//captureName := fmt.Sprintf("capture_%d.%s", CaptureCount, ext)
	//	//_ = filepath.Join(config.Settings.FilePath.Record, template.CurrentTemplate, captureName)
	//	//windowName :=
	//
	//	//x, y, w, h := robotgo.GetBound
	//	//s(win.Pid)
	//	//img := robotgo.CaptureImg(x, y, w, h)
	//	//ImageChan <- imgFile
	//	//if err := robotgo.Save(img, imgFile); err != nil {
	//	//	fmt.Println("==========>", err)
	//	//}
	//	//CaptureCount++
	//}
}

func (c *Capture) CaptureEvent(input chan string) {
	for imgFile := range input {
		//base64Buffer := bytes.NewBuffer(nil)
		//err := jpeg.Encode(base64Buffer, img, nil)
		//if err != nil {
		//	fmt.Println(err)
		//	return
		//}
		////fmt.Println(base64Buffer.String())
		////fmt.Println(len(base64Buffer.Bytes()))
		//dist := make([]byte, 5000000)                         //开辟存储空间
		//base64.StdEncoding.Encode(dist, base64Buffer.Bytes()) //buff转成base64
		// FIXME 这里可能最佳解决方案还是调用window capture 来录屏 周末搞一下

		// 如果不行，那就使用wails的动态资产获取文件路径
		runtime.EventsEmit(c.ctx, "capture", imgFile)
	}
}

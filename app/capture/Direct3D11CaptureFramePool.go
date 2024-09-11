package capture

import (
	"errors"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/lxn/win"
	"github.com/tsqqqqqq/GoRecord/winapi/winrt"
)

// Direct3D11CaptureFramePool

// Protect from gabage collecter
var generatedDirect3D11CaptureFramePool = map[uintptr]*Direct3D11CaptureFramePoolVtbl{}

type Direct3D11CaptureFramePool struct {
	ole.IUnknown
}

type Direct3D11CaptureFramePoolVtbl struct {
	ole.IUnknownVtbl
	Invoke  uintptr
	counter *int
}

// NewDirect3D11CaptureFramePool 解释一下
// invoke 是一个回调函数，这个回调函数的类型是winrt.Direct3D11CaptureFramePoolFrameArrivedProcType
// 这个回调函数的参数是一个指针，这个指针指向的是一个IDirect3D11CaptureFramePool
// 这个回调函数的返回值是一个error
// 这个回调函数的作用是当有新的帧到达的时候，就会调用这个回调函数
// 也就是说，这个回调函数的作用是当有新的帧到达的时候，就会调用这个回调函数
func NewDirect3D11CaptureFramePool(invoke winrt.Direct3D11CaptureFramePoolFrameArrivedProcType) *Direct3D11CaptureFramePool {
	var counter = 1
	var v = &Direct3D11CaptureFramePoolVtbl{
		Invoke:  syscall.NewCallback(invoke),
		counter: &counter,
	}
	var newV = new(Direct3D11CaptureFramePool)
	newV.RawVTable = (*interface{})(unsafe.Pointer(v))

	v.QueryInterface = syscall.NewCallback(newV.queryInterface)
	v.AddRef = syscall.NewCallback(newV.addRef)
	v.Release = syscall.NewCallback(newV.release)

	generatedDirect3D11CaptureFramePool[uintptr(unsafe.Pointer(newV))] = v

	return newV
}

func (v *Direct3D11CaptureFramePool) VTable() *Direct3D11CaptureFramePoolVtbl {
	return (*Direct3D11CaptureFramePoolVtbl)(unsafe.Pointer(v.RawVTable))
}

func (v *Direct3D11CaptureFramePool) Invoke(sender *winrt.IDirect3D11CaptureFramePool, args *ole.IInspectable) error {
	r1, _, _ := syscall.SyscallN(v.VTable().Invoke, uintptr(unsafe.Pointer(sender)), uintptr(unsafe.Pointer(args)))
	return ole.NewError(r1)
}

// QueryInterface(vp *Direct3D11CaptureFramePool, riid ole.GUID, lppvObj **ole.Inspectable)
func (v *Direct3D11CaptureFramePool) queryInterface(lpMyObj *uintptr, riid *uintptr, lppvObj **uintptr) uintptr {
	// Validate input
	if lpMyObj == nil {
		return win.E_INVALIDARG
	}

	var V = new(Direct3D11CaptureFramePool)

	var err error
	// Check dereferencability
	func() {
		defer func() {
			if recover() != nil {
				err = errors.New("InvalidObject")
			}
		}()
		// if object cannot be dereferenced, then panic occurs
		*V = *(*Direct3D11CaptureFramePool)(unsafe.Pointer(lpMyObj))
		V.VTable()
	}()
	if err != nil {
		return win.E_INVALIDARG
	}

	*lppvObj = nil
	var id = new(ole.GUID)
	*id = *(*ole.GUID)(unsafe.Pointer(riid))

	// Convert
	switch id.String() {
	case ole.IID_IUnknown.String(), winrt.ITypedEventHandlerID.String(), winrt.IAgileObjectID.String():
		V.AddRef()
		*lppvObj = (*uintptr)(unsafe.Pointer(V))

		return win.S_OK
	default:
		return win.E_NOINTERFACE
	}
}

func (v *Direct3D11CaptureFramePool) addRef(lpMyObj *uintptr) uintptr {
	// Validate input
	if lpMyObj == nil {
		return 0
	}

	var V = (*Direct3D11CaptureFramePool)(unsafe.Pointer(lpMyObj))
	*V.VTable().counter++

	return uintptr(*V.VTable().counter)
}

func (v *Direct3D11CaptureFramePool) release(lpMyObj *uintptr) uintptr {
	// Validate input
	if lpMyObj == nil {
		return 0
	}

	var V = (*Direct3D11CaptureFramePool)(unsafe.Pointer(lpMyObj))
	*V.VTable().counter--

	if *V.VTable().counter == 0 {
		V.RawVTable = nil
		_, ok := generatedDirect3D11CaptureFramePool[uintptr(unsafe.Pointer(lpMyObj))]
		if ok {
			delete(generatedDirect3D11CaptureFramePool, uintptr(unsafe.Pointer(lpMyObj)))
			runtime.GC()
		}
		return 0
	}

	return uintptr(*V.VTable().counter)
}

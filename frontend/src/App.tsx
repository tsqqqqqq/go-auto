import {SetStateAction, useEffect, useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Button, Flex, Form, Image, Input, Layout, Modal, notification, Select} from "antd";
import {OnListen, Run} from "../wailsjs/go/auto/AutoRecord";
import {ChangeCurrentTemplate, GetAll, CreateTemplate} from "../wailsjs/go/template/Template";
// @ts-ignore
import {ChangeCurrentWindow, GetWindows} from "../wailsjs/go/window/Window";

// @ts-ignore
import {template, window as win} from "../wailsjs/go/models";
import Template = template.Template;
import Window = win.Window

import FormItem from "antd/es/form/FormItem";
import {Content, Footer} from "antd/es/layout/layout";
import {EventsOn} from "../wailsjs/runtime";


type NotificationType = 'success' | 'info' | 'warning' | 'error';

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const [listen, setListen] = useState(false)

    const [templates, setTemplates] = useState<Array<Template>>([])
    const [selectedTemplate, setSelectedTemplate] = useState<string | null>(null)

    const [windows, setWindows] = useState<Array<Window>>([])
    const [selectWindow, setSelectWindow] = useState<string | null>(null)

    const [isModalOpen, setIsModalOpen] = useState(false)


    const [api, contextHolder] = notification.useNotification();

    const openNotificationWithIcon = (type: NotificationType, title: string, text: string) => {
        api[type]({
            message: title,
            description: text
        });
    };

    const templateFields = {
        label: "Name",
        value: "Name"
    }

    const windowFields = {
        label: "Title",
        value: "Pid"
    }

    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);

    const initWindows = () => {
        GetWindows().then((res: SetStateAction<Window[]>) => {
            setWindows(res)
        })
    }

    const initTemplates = () => {
        GetAll().then((res) => {
            setTemplates(res)
        })
    }

    useEffect(() => {
        initWindows()
        initTemplates()
        eventTest()
    }, [])

    useEffect(() => {
        initTemplates()
    }, [isModalOpen]);

    // 不传递的时候默认为isListen的相反数, 当重播的时候会发送停止监听的信号
    const handleListen = (isListen: boolean = !listen) => {
        OnListen(isListen).then(() =>
            setListen(isListen)
        )
    }

    const handleRun = () => {
        //如果这个时候还在录制则停止，否则直接执行
        if (listen) {
            handleListen(false)
        }
        if (!selectedTemplate) {
            openNotificationWithIcon('error', 'please select template', '')
            return
        }
        Run().then(() => {

        })
    }

    const handleTemplateChange = (value:string) => {
        ChangeCurrentTemplate(value).then(() => {
            setSelectedTemplate(value)
        })
    }

    const handleWindowChange = (value:string) => {
        const item: win.Window | any = windows.find((item: win.Window | any) => item['Pid'] === value)
        ChangeCurrentWindow(item)
    }

    const layoutStyle = {
        borderRadius: 8,
        overflow: 'hidden',
        width: 'calc(90% - 8px)',
        maxWidth: 'calc(90% - 8px)',
    };

    const contentStyle: React.CSSProperties = {
        textAlign: 'center',
        minHeight: '50%',
        color: '#fff',
    };

    const eventTest = () => {
        try {
            EventsOn('capture', (count: any) => {
                console.log(count)
            })
        } catch (err) {
            console.log(err)
        }
    }



    return (
        <>
            {contextHolder}
            <Flex className='justify-center'>
                <Layout style={layoutStyle} className='shadow-2xl h-full min-h-full mt-10'>
                    <Content style={contentStyle} className=''>
                        <img src={logo} id="logo" alt="logo"/>

                    </Content>
                    <Footer className={"w-full"}>
                            <Form
                                layout={"horizontal"}
                                labelCol={{ span: 8 }}
                                wrapperCol={{ span: 24 }}
                                labelAlign={"right"}
                                initialValues={{ remember: true }}
                                autoComplete="off"
                            >
                                <FormItem label={"Window: "} className={'w-1/2 m-auto mb-3'}>
                                    <Select onChange={handleWindowChange} fieldNames={windowFields} options={windows} onFocus={initWindows}></Select>
                                </FormItem>

                                <FormItem label={"Templates: "} className={'w-1/2 m-auto mb-3'}>
                                    <Select onChange={handleTemplateChange} fieldNames={templateFields} options={templates}></Select>
                                </FormItem>
                                <FormItem>
                                    <div className={' m-auto'}>
                                        <Button className="btn ml-10" onClick={() => handleListen()}> {listen ? 'Stop': 'Listen'} </Button>
                                        <Button className="btn ml-10" onClick={() => handleRun()}> Run </Button>
                                        <Button className={"btn ml-10"} onClick={() => setIsModalOpen(!isModalOpen)}> New </Button>
                                    </div>
                                </FormItem>
                            </Form>

                    </Footer>
                </Layout>
            </Flex>
            <TemplateConfigModal isModalOpen={isModalOpen} setIsModalOpen={setIsModalOpen} openNotificationWithIcon={openNotificationWithIcon}/>
        </>

    )
}

// @ts-ignore
const TemplateConfigModal = ({isModalOpen, setIsModalOpen, openNotificationWithIcon}) => {
    const [form] = Form.useForm<Template>()

    const handleOK = () => {
        const data = form.getFieldsValue()
        // @ts-ignore
        CreateTemplate(data.name).then((res) => {
            if (res) {
                openNotificationWithIcon('success', 'Create Template success', '')
                setIsModalOpen(false)
            } else {
                console.log(res)
            }

        }).catch((err) => {
            openNotificationWithIcon('error', 'Create Template Fail', err)
        })
    }

    const handleCancel = () => {
        setIsModalOpen(false)
    }

    return (
        <>
            <Modal  title={"Configuration Template"} open={isModalOpen} onOk={handleOK} onCancel={handleCancel}>
                <Form

                    form={form}
                >
                    <FormItem name="name" label={'Name'}>
                        <Input />
                    </FormItem>
                </Form>
            </Modal>
        </>
    )
}

export default App

import {useEffect, useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";
import {Button, Flex, Form, Image, Layout, Select} from "antd";
import {OnListen, Run} from "../wailsjs/go/auto/AutoRecord";
import {ChangeCurrentTemplate, GetAll} from "../wailsjs/go/template/Template";
import {template} from "../wailsjs/go/models";
import Template = template.Template;
import FormItem from "antd/es/form/FormItem";
import {Content, Footer} from "antd/es/layout/layout";


function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const [listen, setListen] = useState(false)

    const [templates, setTemplates] = useState<Array<Template>>([])
    const [selectedTemplate, setSelectedTemplate] = useState('')
    const templateFields = {
        label: "Name",
        value: "Name"
    }

    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);

    const initTemplates = () => {
        GetAll().then((res) => {
            console.log(res)
            setTemplates(res)
        })
    }

    useEffect(() => {
        initTemplates()
    }, [])

    function greet() {
        Greet(name).then(updateResultText);
    }

    // 不传递的时候默认为isListen的相反数, 当重播的时候会发送停止监听的信号
    const handleListen = (isListen: boolean = !listen) => {
        OnListen(isListen).then(() =>
            setListen(isListen)
        )
    }

    const handleRun = () => {
        handleListen(false)
        Run().then(() => {

        })
    }

    const handleTemplateChange = (value:string) => {
        ChangeCurrentTemplate(value).then(() => {
            setSelectedTemplate(value)
        })
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

    const footerStyle: React.CSSProperties = {
        textAlign: 'center',
        color: '#fff',
        // backgroundColor: '#4096ff',
    };

    return (
            <Flex className='justify-center'>
                <Layout style={layoutStyle} className='shadow-2xl h-full min-h-full mt-10'>
                    <Content style={contentStyle} className=''>
                        <img src={logo} id="logo" alt="logo"/>
                    </Content>
                    <Footer style={footerStyle}>
                        <Form
                            layout={"inline"}
                            className={"w-full"}
                            initialValues={{ remember: true }}
                            autoComplete="off"
                        >
                            <FormItem label={"Templates: "} className={'w-1/2'}>
                                <Select onChange={handleTemplateChange} fieldNames={templateFields} options={templates}></Select>
                            </FormItem>
                            <FormItem>
                                <Button className="btn" onClick={() => handleListen()}> {listen ? 'Stop': 'Listen'} </Button>
                            </FormItem>
                            <FormItem>
                                <Button className="btn" onClick={() => handleRun()}> Run </Button>
                            </FormItem>
                        </Form>
                    </Footer>
                </Layout>
            </Flex>
    )
}

export default App

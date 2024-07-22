import {useEffect, useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";
import {Button, Form, Select} from "antd";
import {OnListen, Run} from "../wailsjs/go/auto/AutoRecord";
import {ChangeCurrentTemplate, GetAll} from "../wailsjs/go/template/Template";
import {template} from "../wailsjs/go/models";
import Template = template.Template;
import FormItem from "antd/es/form/FormItem";


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

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
            <div >
                <Form

                    style={{ maxWidth: 200 }}
                    initialValues={{ remember: true }}
                    autoComplete="off"
                >
                    <FormItem label={"Templates: "}>
                        <Select onChange={handleTemplateChange} fieldNames={templateFields} options={templates}></Select>
                    </FormItem>
                </Form>
            </div>

            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <Button className="btn" onClick={greet}>Greet</Button>
                <Button className="btn" onClick={() => handleListen()}> {listen ? 'Stop': 'Listen'} </Button>
                <Button className="btn" onClick={() => handleRun()}> Run </Button>
            </div>
        </div>
    )
}

export default App

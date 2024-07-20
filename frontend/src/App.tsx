import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";
import {Button} from "antd";
import {OnListen, Run} from "../wailsjs/go/auto/AutoRecord";

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const [listen, setListen] = useState(false)

    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);


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

    return (
        <div id="App">
            <img src={logo} id="logo" alt="logo"/>
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

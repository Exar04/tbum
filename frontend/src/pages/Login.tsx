import { EventHandler, useState } from "react"
import { useWebSocketSender, useSocketSubscribe } from "../context/webSocketContext"
import { useNavigate, Link } from "react-router-dom"
import { useAuth } from "../context/authContext"

export const LoginPage = () => {
    const sendthismessagetoserver = useWebSocketSender()
    const [usernameInInput, setUsernameInInput] = useState("")
    const navigate = useNavigate();

    const { setUserId} = useAuth();

    function sendMessage() {
        if (usernameInInput == ""){
            return
        }
        const jsontosend = {
            messageType: "newUser", 
            sender: usernameInInput,
        }
        sendthismessagetoserver(JSON.stringify(jsontosend))
        setUserId(usernameInInput) 
        navigate("/")
    }
    // function hala(da: any){
    //     console.log(da)
    // }
    // useSocketSubscribe(hala)


    return (
        <div className=" flex justify-center items-center w-screen h-screen font-mono font-bold text-xl bg-gradient-to-br from-cyan-400 to-blue-400">
            Login with Username
            <input onChange={(e) => {setUsernameInInput(e.target.value)}} className=" border-2 rounded-lg m-3 outline-none px-2"/>
            <div role={"button"} onClick={() => {sendMessage()}} className=" bg-cyan-400 rounded-full p-2 text-white">Go</div>
        </div>
    )
}
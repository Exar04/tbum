import { FC, useEffect, useState } from "react"
import { useWebSocketSender, useSocketSubscribe } from "../context/webSocketContext"
import { apiMessage, chatMapTypes } from "../types/types"
import { useParams } from 'react-router-dom';
import { useAuth } from "../context/authContext";

interface messagesProps{
    // idOfCurrentChatOpened : string
    setChatMap: React.Dispatch<React.SetStateAction<chatMapTypes[] | undefined>>
    chatMap: chatMapTypes[] | undefined 
}

export const Messages: FC<messagesProps> = ({chatMap, setChatMap}) => {
    const { userid} = useAuth()
    const [textInTextarea, setTextInTextarea] = useState("")
    const { chatId } = useParams<{ chatId: string }>();
    const sendthismessagetoserver = useWebSocketSender()
    function sendMessage() {
        if (textInTextarea == ""){
            return
        }

        const jsontosend = {
            content: textInTextarea,
            sender: userid,
            reciever:chatId,
            messageType: "userMessage", 
        }

        const addtomap: chatMapTypes =  {
            content: textInTextarea,
            sender: userid,
            reciever: chatId ?? "",
        }

        setChatMap([...(chatMap ?? []), addtomap]);

        sendthismessagetoserver(JSON.stringify(jsontosend))
    }

    useEffect(() => {
      console.log(chatMap)
    
    }, [chatMap])
    

    const selectedUserChats = (chatMap??[]).map((chatas:chatMapTypes, index) => (
        <div key={index}>
            {
                chatas.sender == chatId ? 
                <div>
                    {chatas.content}
                </div> : ""
            }
            {
                chatas.sender == userid && chatas.reciever == chatId?  
                <div className=" bg-red-500">
                    {chatas.content}
                </div> : ""
            }
        </div>
    )) 

    return (
        <div className=" flex flex-col relative w-full h-full ">
            <div className=" w-full bg-red-400 h-12 flex p-2 items-center ">
                <div className=" bg-slate-500 w-8 h-8 rounded-full mx-2"></div>
                <div className=" font-mono text-white"> {chatId} </div>
            </div>
            <div className=" bg-emerald-600 flex-grow overflow-y-scroll">
               {selectedUserChats} 
            </div>
            <div className="absolute bottom-0 left-0  w-full bg-red-300 p-2">
                <div className="flex justify-center items-center bg-white rounded-lg p-1">
                <textarea onChange={(e) => {setTextInTextarea(e.target.value)}} className="  outline-none w-full h-full rounded-lg p-2"/>
                <div role={"button"} onClick={() => {sendMessage()}} className=" bg-cyan-400 rounded-full p-2 font-mono text-white duration-100 hover:scale-125 hover:-translate-x-2">Send</div>
                </div>
            </div>
        </div>
    )
}
import React, { useState } from "react"
import { Outlet, useParams } from "react-router-dom"
import { Messages } from "../components/MessagesPage"
import { Sidebar } from "../components/Sidebar"
import { chatMapTypes } from "../types/types"

interface ChatPageProps {
    // setUserId: React.Dispatch<React.SetStateAction<number | undefined>>
    // setLogIn: React.Dispatch<React.SetStateAction<boolean>>
    // setUserName: React.Dispatch<React.SetStateAction<string>>
    // receivedData: ApiMessage | undefined
}

export const ChatPage: React.FC<ChatPageProps> = ({}) => {
    // const [idOfCurrentChatOpened, setIdOfCurrentChatOpened] = useState("")
    const { chatId } = useParams<{ chatId: string }>();
    const [chatMap, setChatMap] = useState<chatMapTypes[]>()
    return(
        <div className=" h-screen w-screen flex">
            <Sidebar setChatMap={setChatMap} chatMap={chatMap}/>
            <div className="no-scrollbar w-full h-full">
                {/* <Outlet /> */}
                {chatId ?<Messages setChatMap={setChatMap} chatMap={chatMap}/>: <div className=" w-full h-full font-mono font-bold flex justify-center items-center text-2xl"> No chat selected</div>}
            </div>
        </div>
    )
}
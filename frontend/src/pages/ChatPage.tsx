import React from "react"
import { Outlet } from "react-router-dom"
import { Sidebar } from "../components/Sidebar"

interface ChatPageProps {
    // setUserId: React.Dispatch<React.SetStateAction<number | undefined>>
    // setLogIn: React.Dispatch<React.SetStateAction<boolean>>
    // setUserName: React.Dispatch<React.SetStateAction<string>>
    // receivedData: ApiMessage | undefined
}

export const ChatPage: React.FC<ChatPageProps> = ({}) => {
    return(
        <div className=" h-screen w-screen flex">
            <Sidebar />
            <div className="no-scrollbar w-full h-full">
                <Outlet />
            </div>
        </div>
    )
}
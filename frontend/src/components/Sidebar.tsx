import { FC, useState } from "react"
import { useNavigate } from "react-router-dom"
import { useSocketSubscribe } from "../context/webSocketContext"
import { apiMessage, chatMapTypes } from "../types/types"

interface sidebarProps {
    // setIdOfCurrentChatOpened: React.Dispatch<React.SetStateAction<string>>
    setChatMap: React.Dispatch<React.SetStateAction<chatMapTypes[] | undefined>>
    chatMap: chatMapTypes[] | undefined 
}

export const Sidebar: FC<sidebarProps> = ({setChatMap, chatMap}) => {
    const [userIdfromInput, setUserIdFromInput] = useState("")
    const [friendslist, setFriendslist] = useState<string[]>([])

    const handleSocketUpdate = (receivedData: apiMessage) => {
        console.log(receivedData.content)
        const addtomap: chatMapTypes =  {
            content: receivedData.content,
            sender: receivedData.sender,
            reciever: receivedData.reciever
        }
        setChatMap([...(chatMap ?? []), addtomap]);
    }

    useSocketSubscribe(handleSocketUpdate)

    return (
        <div className=" flex-none w-2/6 h-full bg-slate-400 p-3">
            <div className=" text-white font-mono text-3xl font-bold ">Chats</div>
            <div className=" relative flex">
                {/* <img role={"button"} src="icons/black-add-symbol.png" className=" w-8 h-8 absolute top-3 right-1 bg-slate-400 rounded-full p-1"/> */}
                <div role={"button"} onClick={() => {setFriendslist([...friendslist,userIdfromInput])}} className=" w-fit h-8 absolute top-3 right-1 bg-slate-400 rounded-full px-2 flex justify-center items-center hover:scale-110 duration-100">Add</div>
                <input onChange={(e) => {setUserIdFromInput(e.target.value)}} className={` rounded-xl w-full h-10 p-2 outline-none font-mono my-2`} />
            </div>
            <FriendsList friendslist={friendslist}/>
        </div>
    )
}

interface friendslistProps{
    // setIdOfCurrentChatOpened: React.Dispatch<React.SetStateAction<string>>
    friendslist: string[]
}

const FriendsList: FC<friendslistProps> = ({friendslist}) => {
    const navi = useNavigate()

    const fList = friendslist.map((user, index) => (
        <div key={index} onClick={() => {navi(`/messages/${user}`)}} className=" p-2 flex items-center h-20 hover:translate-x-8 duration-75 border-b-0.5 ">
            <div className=" w-10 h-10 bg-slate-300 rounded-full mx-2 flex-none"></div>
            <div className=" h-full w-full font-mono text-white p-2">
                <div>{user}</div>
                <div className=" text-xs">Message</div> 
            </div>
        </div>
    ))
    return (
        <div>
            {fList}
        </div>
    )
}
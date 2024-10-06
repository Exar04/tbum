import { FC, useState } from "react"

interface sidebarProps{

}

export const Sidebar: FC<sidebarProps> = () => {
    return (
        <div className=" flex-none w-2/6 h-full bg-slate-400 p-3">
            <div className=" text-white font-mono text-3xl font-bold ">Chats</div>
            <input className={` rounded-xl w-full h-10 p-2 outline-none font-mono my-2`} />
            <FriendsList />
        </div>
    )
}

interface friendslistProps{}
const FriendsList: FC<friendslistProps> = () => {
    const a =  ["la", "ba", "ka"]
    const fList = a.map((users) => (
        <div className=" p-2 flex items-center h-20 hover:translate-x-8 duration-75 border-b-0.5 ">
            <div className=" w-10 h-10 bg-slate-300 rounded-full mx-2 flex-none"></div>
            <div className=" h-full w-full font-mono text-white p-2">
                <div>{a}</div>
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
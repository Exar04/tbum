import { FC } from "react"

interface messagesProps{

}

export const Messages: FC<messagesProps> = () => {
    return (
        <div className=" flex flex-col relative w-full h-full bg-green-300">
            <div className=" w-full bg-red-400 h-10 flex p-2 items-center ">
                <div className=" bg-slate-500 w-8 h-8 rounded-full mx-2"></div>
                <div className=" font-mono text-white"> Username </div>
            </div>
            <div className="absolute bottom-0 left-0  w-full bg-red-300 p-2">
                <div className="flex justify-center items-center bg-white rounded-lg p-1">
                <textarea className="  outline-none w-full h-full rounded-lg p-2"/>
                <div role={"button"} className=" bg-cyan-400 rounded-full p-2 font-mono text-white duration-100 hover:scale-125 hover:-translate-x-2">Send</div>
                </div>
            </div>
        </div>
    )
}
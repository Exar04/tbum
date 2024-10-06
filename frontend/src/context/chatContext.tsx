import { useContext, createContext, ReactNode } from "react"

interface ChatContextType {
    message: string
    sendMessage: (msg: string) => void
}

const ChatContext  = createContext < ChatContextType | undefined >(undefined)

export function useChat() {
    return useContext(ChatContext)
}

export function ChatProvider({ children  }:{ children: ReactNode }) { 

    const sendMessage = (msg: string) => {
        console.log('Message sent:', msg);
    };

    const value : ChatContextType = {
        message: "",
        sendMessage
    }

    return(
        <ChatContext.Provider value={value}>
            {children}
        </ChatContext.Provider>
    )
}
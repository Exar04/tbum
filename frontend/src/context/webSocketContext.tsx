// import { useContext, createContext, ReactNode } from "react"

import { createContext, ReactNode, useContext, useEffect, useRef, useState } from "react";

interface ChatContextType {
    socket: WebSocket | null;
    // message: string
    // sendMessage: (msg: string) => void
}

export const SocketContext = createContext< ChatContextType | undefined >(undefined );

const WebSocketProvider = ({ children }:{ children: ReactNode }) => {
    const socketRef = useRef< WebSocket | null >(null); // Initialize with null
    const [socket, setSocket] = useState<WebSocket | null>(null); // State to hold socket


    useEffect(() => {
        // Create the WebSocket instance
        socketRef.current = new WebSocket('ws://localhost:8088/websocket');

        // Event handlers
        socketRef.current.onopen = () => {
            console.log('WebSocket connection established.');
            setSocket(socketRef.current)
        };

        socketRef.current.onclose = () => {
            console.log('WebSocket connection closed.');
            setSocket(null)
        };

        // Cleanup on unmount
        return () => {
            if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
                socketRef.current.close();
            }
        };
    }, []);

    // const message = ""
    // const sendMessage = (msg: string) => {
    //     if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
    //         socketRef.current.send(msg);
    //     }
    // };
    const contextValue: ChatContextType = {
        // socket: socketRef.current,
        socket: socket,
        // other properties if needed
    };

    return (
        <SocketContext.Provider value={contextValue}>
            {children}
        </SocketContext.Provider>
    );
};

export default WebSocketProvider;

type EventHandler = (message: any) => void;
export const useSocketSubscribe = (eventHandler: EventHandler) => {
    const context = useContext(SocketContext);
    const socket = context?.socket;

    useEffect(() => {
        if (socket) {
            const messageListener = (event: MessageEvent) => {
                const message = JSON.parse(event.data);
                console.log("tf")
                eventHandler(message);
            };
      
            socket.addEventListener('message', messageListener);
      
            return () => {
                console.log('WebSocket: removing message listener');
                socket.removeEventListener('message', messageListener);
            };
        }
    }, [socket, eventHandler]);
};

export const useWebSocketSender = () => {
    const context = useContext(SocketContext);
    console.log("socket context :", context)
    const socket = context?.socket;

    const sendMessage = (message : string) => {
        console.log(message)
        console.log(socket)
        if (socket) {
            console.log('WebSocket connection state:', socket.readyState);
            
            if (socket.readyState === WebSocket.OPEN) {
                socket.send(message);
                console.log('Message sent:', message);
            } else {
                console.log('WebSocket connection is not open.');
            }
        } else {
            console.log('Socket is not available.');
        }
    };

    return sendMessage;
};
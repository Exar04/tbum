import { createContext, ReactNode, useContext, useState } from "react"

interface AuthContextType {
    jwtToken: string
    userid: string
    setUserId: (uid: string) => void
    // getAuth: (msg: string) => void
}

const AuthContext = createContext <AuthContextType | undefined>(undefined)

export function useAuth() {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
}


export function AuthProvider({ children }:{ children: ReactNode }) {
    const [userId, setUserId] = useState<string >("");

    const value: AuthContextType = {
        jwtToken: "",
        userid:userId,
        setUserId: setUserId 
    }
      
    return(
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    )
}
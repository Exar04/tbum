import { createContext, ReactNode, useContext } from "react"

interface AuthContextType {
    jwtToken: string
    // getAuth: (msg: string) => void
}

const AuthContext = createContext <AuthContextType | undefined>(undefined)

export function useAuth() {
    return useContext(AuthContext)
}

export function AuthProvider({ children }:{ children: ReactNode }) {

    
    const value: AuthContextType = {
        jwtToken: "",
      }
      
    return(
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    )
}
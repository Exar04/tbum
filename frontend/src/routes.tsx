import { createBrowserRouter } from "react-router-dom";
import React, { useState, ComponentType } from "react";
import { Navigate } from "react-router-dom";

import { useAuth } from "./context/authContext";
import { ChatPage } from "./pages/ChatPage";
import { NotFound } from "./pages/NotFound";
import { Profile } from "./components/ProfilePage";
import { Messages } from "./components/MessagesPage";
import { EditProfile } from "./components/EditProfilePage";
import { LoginPage } from "./pages/Login";
import { SignupPage } from "./pages/Signup";


interface PrivateRouteProps {
    component: React.ComponentType<any>; 
    // authenticated: boolean
  }

export const PrivateRoute: React.FC<PrivateRouteProps> = ({
  component: Component,
//   authenticated,
  ...rest
}) => {
//   const { currentUser } = useAuth();
  return true ? (
    <Component {...rest} />
  ) : (
    <Navigate to="/login" replace />
  );
};

export const router = createBrowserRouter([
  {
    path: "/login",
    element: <LoginPage />,
  },
  {
    path: "/signup",
    element: <SignupPage />,
  },
  {
    path: "/",
    element: <PrivateRoute component={ChatPage} />,
    errorElement: <NotFound />,
    children: [
      {
        path: "/profile/:profileId",
        element: <Profile />,
      },
      {
        path: "/messages/:chatId",
        element: <Messages />,
      },
      {
        path: "/editprofile",
        element: <EditProfile />,
      },
    ],
  },
]);
import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';

import { router } from './routes';
import { RouterProvider } from 'react-router-dom';
import {AuthProvider} from './context/authContext'
import { ChatProvider } from './context/chatContext';


const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
)
root.render(
  <React.StrictMode>
    {/* <AuthProvider> */}
      {/* <ChatProvider> */}
        <RouterProvider router={router}/>
      {/* </ChatProvider> */}
    {/* </AuthProvider> */}
  </React.StrictMode>
)

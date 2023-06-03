"use client";
import { ToastContainer } from "react-toastify";

const NotistackProvider = () => (
  <ToastContainer
    position="top-right"
    autoClose={false}
    hideProgressBar={false}
    newestOnTop={false}
    closeOnClick
    rtl={false}
    pauseOnFocusLoss
    draggable
    pauseOnHover
    theme="light"
  />
);

export default NotistackProvider;

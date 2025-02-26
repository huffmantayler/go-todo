import React from "react";
import { Alert as MuiAlert, Snackbar } from "@mui/material";

const Alert = ({ open, message, status, onClose }) => {
    return (
    <Snackbar open={open} autoHideDuration={6000} onClose={onClose}>
      <MuiAlert onClose={onClose} severity={status} sx={{ width: "100%"}}>
        {message}
      </MuiAlert>
    </Snackbar>
  );
};

export default Alert;

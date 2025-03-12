import React, { useState } from "react";
import { TextField, Button, List, ListItem, ListItemText, IconButton, Paper, Typography } from "@mui/material";
import Alert from "./Alert";
import { Delete } from "@mui/icons-material";
import { useEffect } from "react";
import axios from 'axios';

export default function TodoList() {
  const [todos, setTodos] = useState([]);
  const [input, setInput] = useState("");
  const backendUrl = import.meta.env.VITE_BACKEND_URL;
  const [alert, setAlert] = useState({ open: false, message: "", status: "" });
  const [editing, setEditing] = useState({editing: false, id: null})
  const [todoTitle, setTodoTitle] = useState("")


  const handleApiResponse = (status, message) => {
    setAlert({
      open: true,
      message: message,
      status: status,
    });
  };

  const handleCloseAlert = () => {
    setAlert({ open: false, message: "", status: "" });
  };

  useEffect(() =>{
    getAllTodos()
  }, [])

  const getAllTodos = () => {
    axios.get(`${backendUrl}/getAllTodos`)
    .then(res => {
      console.log("Frontend Changes1") 
      setTodos(res.data)
    })
    .catch(err => {
        handleApiResponse("error", err.message)
        console.log(err)
    })
  }

  const addTodo = () => {
    if (!input.trim()) return;
    createTodo(input)
    setTodos([...todos, { id: Date.now(), title: input, done: false }]);
    setInput("");
  };

  const toggleTodo = (id, done) => {

    axios.post(`${backendUrl}/updateTodo?&id=${id}&done=${!done}`)
    .then(res =>
        handleApiResponse('success', res.data)
    ).catch(err => {
        handleApiResponse('error', err)
    })

    setTodos(todos.map(todo => 
      todo.id === id ? { ...todo, done: !todo.done } : todo
    ));
  };

  const updateTodoTitle = (id) => {
    axios.post(`${backendUrl}/updateTodo?&id=${id}&title=${todoTitle}`)
    .then(res => {
        handleApiResponse('success', res.data)
    }).catch(err => {
        handleApiResponse('error', err)
    })
    setTodos(todos.map(todo => 
        todo.id === id ? { ...todo, title: todoTitle} : todo
    ));
    setTodoTitle("")
    setEditing({editing: false, id: null})
  }

  const deleteTodo = (id) => {
    axios.delete(`${backendUrl}/deleteTodo?id=${id}`)
    .then(res => {
       handleApiResponse("success", res.data)
    })
    .catch(err => {
        handleApiResponse("Error", err)
    })
    setTodos(todos.filter(todo => todo.id !== id));
  };

  const createTodo = (title) => {
    axios.post(`${backendUrl}/createTodo`, {title: title})
    .then(res =>{
        // setAlert({open: true, message: res.data, status: "success"})
        handleApiResponse("success", res.data)
    })
    .catch(err => {
        handleApiResponse("error", err)
    })
  }

  const getTaskEmoji = (done) => {
    return done ? "✅" : "🔲";
  };

  return (
    <Paper sx={{ maxHeight: "70vh", maxWidth: 400, margin: "auto", mt: 4, p: 3, textAlign: "center", overflow: "scroll" }}>
    
    {alert.open && <Alert
        open={alert.open}
        message={alert.message}
        status={alert.status}
        onClose={handleCloseAlert}
    />}
      <Typography variant="h5" gutterBottom>Todo List</Typography>

      {/* Input Field */}
      <TextField
        fullWidth
        variant="outlined"
        placeholder="Add a task..."
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={(e) => e.key === "Enter" && addTodo()}
      />
      <Button 
        className="addTaskButton"
        sx={{ mt: 2 }} 
        variant="contained" 
        onClick={addTodo} 
        fullWidth
      >
        Add Task
      </Button>

      {/* Todo Items */}
      <List sx={{ mt: 2, maxHeight: "70vh" }}>
        {todos && todos.map(todo => (
          <ListItem key={todo.id} sx={{ display: "flex", alignItems: "center" }}>
            {
            <div
                onClick={() => toggleTodo(todo.id, todo.done)}
            >{getTaskEmoji(todo.done)}</div>
            }
            {editing && todo.id === editing.id ? 
            <>
            <TextField defaultValue={todo.title} data-testid="task-item" onChange={(e) => setTodoTitle(e.target.value)}/>
            <Button style={{marginLeft: "10px"}} variant="contained" onClick={() => updateTodoTitle(todo.id)}>Save</Button>
            </>
            : 
            <>
                <ListItemText
                className="deleteTaskButton"
                primary={todo.title}
                sx={{ textDecoration: todo.done ? "line-through" : "none" }}
                onClick={() => setEditing({editing: true, id: todo.id})}
                />
                <IconButton onClick={() => deleteTodo(todo.id)} color="error">
                <Delete />
                </IconButton>
            </>
            
        }
          </ListItem>
        ))}
      </List>
    </Paper>
  );
}

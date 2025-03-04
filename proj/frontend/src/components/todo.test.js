import React from "react";
import {screen, fireEvent, waitFor, render} from '@testing-library/react';
import TodoList from './TodoList'

global.importMeta = { env: { VITE_BACKEND_URL: "http://localhost:8080" } };


it('Task is added on click', async () => {

    render(<TodoList/>);

    const inputElement = screen.getByPlaceholderText("Add a task...");
    const addButton = screen.getByText("Add Task")

    fireEvent.change(inputElement, {target: {value: "New Task"}});

    fireEvent.click(addButton);
    
    await waitFor(() => {
        expect(screen.queryByText("New Task")).not.toBeNull();
    });

    
})

'use client';

import React, { useState, useEffect } from 'react';
import { TodoItem } from './TodoItem';
import { TodoForm } from './TodoForm';
import {api, Todo} from "../lib/api";
import { wsClient } from '../lib/websocket';

/**
 * TodoList component manages the display and interaction with todo items.
 * Handles real-time updates through WebSocket connection and provides
 * CRUD operations for todo items.
 */
export function TodoList() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  /**
   * Fetches all todos from the API and updates the local state.
   * Handles loading and error states during the fetch operation.
   */
  const fetchTodos = async () => {
    try {
      const data = await api.getTodos();
      setTodos(data);
      setError(null);
    } catch (err) {
      setError('Failed to load todos');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let mounted = true;

    /**
     * Sets up WebSocket connection and event handlers for real-time todo updates.
     * Establishes connection to the WebSocket server and configures handlers for
     * todo creation, updates, and deletions.
     */
    const setupWebSocket = () => {
          if (!mounted) return;
    
          try {
            if (!wsClient) {
              console.error('WebSocket client is not initialized');
              return;
            }
    
            // Configure handler for newly created todos
            wsClient.onTodoCreated = (todo) => {
              if (!todo || !mounted) return;
              setTodos(prev => [todo, ...prev]);
            };
            
            // Configure handler for updated todos
            wsClient.onTodoUpdated = (todo) => {
              if (!todo || !mounted) return;
              setTodos(prev => prev.map(t => t.id === todo.id ? todo : t));
            };
            
            // Configure handler for deleted todos
            wsClient.onTodoDeleted = (todoId) => {
              if (!todoId || !mounted) return;
              setTodos(prev => prev.filter(t => t.id !== todoId));
            };
    
            // Establish WebSocket connection
            wsClient.connect();
            setIsConnected(true);
          } catch (error) {
            console.error('Error setting up WebSocket handlers:', error);
            setError('Failed to connect to real-time updates');
          }
    
          if (!mounted) {
            console.log('WebSocket setup aborted: component unmounted');
          }
        };

    fetchTodos();
    setupWebSocket();
    
    return () => {
      mounted = false;
      if (wsClient) {
        try {
          wsClient.disconnect();
          setIsConnected(false);
        } catch (error) {
          console.error('Error disconnecting WebSocket:', error);
        }
      }
    };
  }, []);

  /**
   * Handles adding a new todo item to the list.
   * Updates the local state with the new todo.
   * @param todo - The new todo item to add
   */
  const handleAddTodo = async (todo: Todo) => {
    setTodos(prev => [todo, ...prev]);
  };

  /**
   * Updates the status of a todo item.
   * Sends update to the API and refreshes the todo list.
   * @param id - The ID of the todo to update
   * @param status - The new status ('pending' or 'completed')
   */
  const handleUpdateTodo = async (id: number, status: 'pending' | 'completed') => {
    try {
      await api.updateTodoStatus(id, status);
      fetchTodos();
    } catch (err) {
      setError('Failed to update todo');
    }
  };

  /**
   * Deletes a todo item from the list.
   * Sends delete request to the API and updates local state.
   * @param id - The ID of the todo to delete
   */
  const handleDeleteTodo = async (id: number) => {
    try {
      await api.deleteTodo(id);
      setTodos(prev => prev.filter(todo => todo.id !== id));
    } catch (err) {
      setError('Failed to delete todo');
    }
  };

  if (loading) {
    return <div className="flex justify-center p-4">Loading...</div>;
  }

  return (
    <div className="container mx-auto p-4 max-w-4xl">
      <h1 className="text-3xl font-bold mb-8">My Todos</h1>
      
      <TodoForm onAddTodo={handleAddTodo} />

      {error && (
        <div className="bg-red-100 text-red-700 p-4 rounded mb-4">
          {error}
        </div>
      )}

      <div className="space-y-4">
        {todos.map((todo) => (
          <TodoItem
            key={todo.id}
            todo={todo}
            onStatusChange={handleUpdateTodo}
            onDelete={handleDeleteTodo}
          />
        ))}
        {todos.length === 0 && (
          <div className="text-center text-gray-500 p-8">
            No todos yet. Add one above!
          </div>
        )}
      </div>
    </div>
  );
}
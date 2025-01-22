from typing import Any, Optional
import httpx
from mcp.server.fastmcp import FastMCP

# Initialize FastMCP server
mcp = FastMCP("secretary")

# Constants
API_BASE = "http://localhost:8080/api"
USER_AGENT = "secretary-app/1.0"

async def make_api_request(url: str, method: str = "GET", data: Optional[dict] = None) -> dict[str, Any] | None:
    """Make a request to the TODO API with proper error handling."""
    headers = {
        "User-Agent": USER_AGENT,
        "Accept": "application/json",
        "Content-Type": "application/json"
    }
    
    async with httpx.AsyncClient() as client:
        try:
            if method == "GET":
                response = await client.get(url, headers=headers, timeout=10.0)
            elif method == "POST":
                response = await client.post(url, headers=headers, json=data, timeout=10.0)
            elif method == "PUT":
                response = await client.put(url, headers=headers, json=data, timeout=10.0)
            elif method == "DELETE":
                response = await client.delete(url, headers=headers, timeout=10.0)
            
            response.raise_for_status()
            return response.json() if response.text else None
        except Exception:
            return None

def format_todo(todo: dict) -> str:
    """Format a todo into a readable string."""
    return f"""
Title: {todo.get('title', 'Unknown')}
Description: {todo.get('description', 'No description')}
Status: {todo.get('status', 'Unknown')}
Due Date: {todo.get('due_date', 'No due date')}
"""

@mcp.tool()
async def list_tasks() -> str:
    """List all tasks in the todo list."""
    data = await make_api_request(f"{API_BASE}/todos")
    if not data:
        return "Unable to fetch tasks."
    if not data:
        return "No tasks found."
    todos = [format_todo(todo) for todo in data]
    return "\n---\n".join(todos)

@mcp.tool()
async def create_task(title: str, description: str, due_date: str) -> str:
    """Create a new task in the todo list.
    Args:
        title: Title of the task
        description: Description of task
        due_date: Due date in YYYY-MM-DD format
    """
    data = {
        "title": title,
        "description": description,
        "due_date": due_date,
        "status": "pending"
    }
    
    response = await make_api_request(f"{API_BASE}/todos", method="POST", data=data)
    if not response:
        return "Failed to create task."
    return f"Task created successfully:\n{format_todo(response)}"

@mcp.tool()
async def update_task_status(task_id: int, status: str) -> str:
    """Update the status of a task.
    Args:
        task_id: ID of the task
        status: New status (pending/completed)
    """
    data = {"status": status}
    response = await make_api_request(f"{API_BASE}/todos/{task_id}", method="PUT", data=data)
    if not response:
        return f"Failed to update task {task_id}."
    return f"Task updated successfully:\n{format_todo(response)}"

if __name__ == "__main__":
    # Initialize and run the server
    mcp.run(transport='stdio')
"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

type Task = {
  id: string;
  payload: string;
  status: string;
  agentId?: string;
  attempt: number;
  version: number;
  fencingToken: number;
  leaseExpiresAt?: string;
  createdAt: string;
  updatedAt: string;
};

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [payload, setPayload] = useState("");
  const [creating, setCreating] = useState(false);

  useEffect(() => {
    fetch("http://localhost:8080/api/tasks")
      .then((response) => response.json())
      .then((data) => {
        setTasks(data);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Failed to load tasks:", error);
        setLoading(false);
      });
  }, []);

  const createTask = async () => {
    if (!payload.trim()) {
      return;
    }

    setCreating(true);

    try {
      const response = await fetch("http://localhost:8080/api/tasks", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          payload: payload.trim(),
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to create task");
      }

      setPayload("");

      const tasksResponse = await fetch("http://localhost:8080/api/tasks");

      if (!tasksResponse.ok) {
        throw new Error("Failed to refresh tasks");
      }

      const tasks = await tasksResponse.json();
      setTasks(tasks);
    } catch (error) {
      console.error("Failed to create task:", error);
    } finally {
      setCreating(false);
    }
  };

  return (
    <main className="container">
      <header className="header">
        <div>
          <h1>Task Dashboard</h1>
          <p>Go Agent Control Plane</p>
        </div>
      </header>

      <section className="create-task">
        <input
          type="text"
          value={payload}
          onChange={(event) => setPayload(event.target.value)}
          placeholder="Enter task payload..."
          disabled={creating}
          onKeyDown={(event) => {
            if (event.key === "Enter") {
              createTask();
            }
          }}
        />

        <button
          onClick={createTask}
          disabled={creating || !payload.trim()}
        >
          {creating ? "Creating..." : "Create Task"}
        </button>
      </section>

      <section className="card">
        <div className="card-header">
          <h2>Tasks</h2>
          <span>{tasks.length} tasks</span>
        </div>

        {loading ? (
          <p className="message">Loading tasks...</p>
        ) : tasks.length === 0 ? (
          <p className="message">No tasks found.</p>
        ) : (
          <div className="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Payload</th>
                  <th>Status</th>
                  <th>Agent</th>
                  <th>Version</th>
                  <th>Created</th>
                </tr>
              </thead>

              <tbody>
                {tasks.map((task) => (
                  <tr key={task.id}>
                    <td className="task-id">
                      <Link href={`/tasks/${task.id}`}>{task.id}</Link>
                    </td>
                    <td>{task.payload}</td>
                    <td>
                      <span
                        className={`status status-${task.status.toLowerCase()}`}
                      >
                        {task.status}
                      </span>
                    </td>
                    <td>{task.agentId ?? "-"}</td>
                    <td>{task.version}</td>
                    <td>
                      {new Date(task.createdAt).toLocaleString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </main>
  );
}
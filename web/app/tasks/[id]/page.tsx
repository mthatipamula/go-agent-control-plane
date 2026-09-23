"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";

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

export default function TaskDetailsPage() {
  const params = useParams();
  const taskId = params.id as string;

  const [task, setTask] = useState<Task | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    const loadTask = async () => {
      try {
        const response = await fetch(
          `http://localhost:8080/api/tasks/${taskId}`
        );

        if (!response.ok) {
          throw new Error("Task not found");
        }

        const data = await response.json();
        setTask(data);
      } catch (error) {
        console.error("Failed to load task:", error);
        setError("Failed to load task");
      } finally {
        setLoading(false);
      }
    };

    loadTask();
  }, [taskId]);

  if (loading) {
    return (
      <main className="container">
        <p className="message">Loading task...</p>
      </main>
    );
  }

  if (error || !task) {
    return (
      <main className="container">
        <p className="message">{error || "Task not found"}</p>
        <Link href="/">← Back to tasks</Link>
      </main>
    );
  }

  return (
    <main className="container">
      <Link href="/" className="back-link">
        ← Back to tasks
      </Link>

      <header className="header">
        <h1>Task Details</h1>
        <p>{task.id}</p>
      </header>

      <section className="details-card">
        <div className="detail-row">
          <span>ID</span>
          <strong>{task.id}</strong>
        </div>

        <div className="detail-row">
          <span>Payload</span>
          <strong>{task.payload}</strong>
        </div>

        <div className="detail-row">
          <span>Status</span>
          <span className={`status status-${task.status.toLowerCase()}`}>
            {task.status}
          </span>
        </div>

        <div className="detail-row">
          <span>Agent</span>
          <strong>{task.agentId ?? "-"}</strong>
        </div>

        <div className="detail-row">
          <span>Attempt</span>
          <strong>{task.attempt}</strong>
        </div>

        <div className="detail-row">
          <span>Version</span>
          <strong>{task.version}</strong>
        </div>

        <div className="detail-row">
          <span>Fencing Token</span>
          <strong>{task.fencingToken}</strong>
        </div>

        <div className="detail-row">
          <span>Created</span>
          <strong>{new Date(task.createdAt).toLocaleString()}</strong>
        </div>

        <div className="detail-row">
          <span>Updated</span>
          <strong>{new Date(task.updatedAt).toLocaleString()}</strong>
        </div>
      </section>
    </main>
  );
}
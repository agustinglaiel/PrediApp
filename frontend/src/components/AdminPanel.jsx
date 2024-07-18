import React, { useEffect, useState } from "react";
import { getUsers, deleteUserById } from "../api/users";
import "../styles/AdminPanel.css";

const AdminPanel = () => {
  const [users, setUsers] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      const users = await getUsers();
      setUsers(users);
    } catch (error) {
      setError("Error fetching users. Please try again.");
    }
  };

  const handleDelete = async (id) => {
    try {
      await deleteUserById(id);
      fetchUsers();
    } catch (error) {
      setError("Error deleting user. Please try again.");
    }
  };

  return (
    <div className="admin-container">
      <h2>Registered Users</h2>
      {error && <p className="error">{error}</p>}
      {users.length === 0 ? (
        <div className="no-users-container">
          <div className="no-users">
            Ups! Ningún usuario registrado todavía.
          </div>
        </div>
      ) : (
        <div className="user-grid">
          {users.map((user) => (
            <div key={user._id} className="user-card">
              <p>
                <strong>ID:</strong> {user._id}
              </p>
              <p>
                <strong>First Name:</strong> {user.first_name}
              </p>
              <p>
                <strong>Last Name:</strong> {user.last_name}
              </p>
              <p>
                <strong>Username:</strong> {user.username}
              </p>
              <p>
                <strong>Email:</strong> {user.email}
              </p>
              <div className="button-container">
                <button
                  className="button delete-button"
                  onClick={() => handleDelete(user._id)}
                >
                  Eliminar
                </button>
                <button className="button update-button">Actualizar</button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default AdminPanel;

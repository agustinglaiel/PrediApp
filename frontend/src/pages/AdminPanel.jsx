import React, { useEffect, useState } from "react";
import { getUsers, deleteUserById } from "../api/users";
import { useNavigate } from "react-router-dom";
import "../styles/AdminPanel.css";

const AdminPanel = () => {
  const [users, setUsers] = useState([]);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      const users = await getUsers();
      console.log("Usuarios recibidos:", users); // Ver los usuarios recibidos
      setUsers(users);
    } catch (error) {
      console.error("Error fetching users:", error);
      setError("Error fetching users. Please try again.");
    }
  };

  const handleDelete = async (id) => {
    console.log("Intentando eliminar usuario con ID:", id); // Verifica que el ID sea correcto
    try {
      await deleteUserById(id);
      console.log("Usuario eliminado exitosamente");
      fetchUsers(); // Actualizar la lista de usuarios
    } catch (error) {
      console.error("Error al eliminar usuario:", error);
      setError("Error deleting user. Please try again.");
    }
  };

  const handleUpdate = (id) => {
    navigate(`/update-user/${id}`);
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
          {console.log("Rendering users: ", users)}
          {users.map((user) => (
            <div key={user._id} className="user-card">
              <p>
                <strong>ID:</strong> {user.id}
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
                  onClick={() => handleDelete(user.id)}
                >
                  Eliminar
                </button>
                <button
                  className="button update-button"
                  onClick={() => handleUpdate(user.id)}
                >
                  Actualizar
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default AdminPanel;

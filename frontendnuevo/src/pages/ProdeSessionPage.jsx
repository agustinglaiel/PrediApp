import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import SessionHeader from "../components/pronosticos/SessionHeader";
import Top3FormHeader from "../components/pronosticos/Top3FormHeader";
import DriverSelect from "../components/pronosticos/DriverSelect";
import SubmitButton from "../components/pronosticos/SubmitButton";
import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import WarningModal from "../components/pronosticos/WarningModal"; // Asegúrate de importar desde la carpeta correcta

// Simulamos isRaceSession desde el frontend para consistencia con el backend
const isRaceSession = (sessionName, sessionType) => {
  return sessionName === "Race" && sessionType === "Race";
};

const ProdeSessionPage = () => {
  const { session_id } = useParams(); // Obtenemos session_id de la URL
  const navigate = useNavigate();
  const [formData, setFormData] = useState({
    P1: null,
    P2: null,
    P3: null,
  });
  const [sessionDetails, setSessionDetails] = useState({
    countryName: "Hungary", // Datos dummy por ahora
    flagUrl: "/images/flags/hungary.jpg",
    sessionType: "Qualifying", // Asegúrate de que sessionType sea correcto
    sessionName: "Qualifying", // Añadimos sessionName como prop
    dateStart: "2025-12-02T04:00:00-03:00", // Datos dummy para date_start (ajusta según necesites)
  });
  const [isFormComplete, setIsFormComplete] = useState(false);
  const [showWarningModal, setShowWarningModal] = useState(false);

  // Simulamos un efecto para cargar datos y validar el tiempo (con logs para depuración)
  useEffect(() => {
    // Aquí normalmente harías una solicitud al backend para obtener sessionDetails, pero por ahora usamos datos dummy
    console.log(`Cargando datos para sesión ${session_id}`);

    // Validar si faltan menos de 5 minutos para dateStart con logs detallados
    const checkTimeRemaining = () => {
      const now = new Date();
      const sessionStart = new Date(sessionDetails.dateStart);

      // Logs detallados de las fechas
      console.log("Fecha y hora actual (now):", now.toISOString());
      console.log(
        "Fecha y hora de inicio de la sesión (dateStart):",
        sessionDetails.dateStart
      );
      console.log(
        "Fecha y hora de inicio de la sesión (parsed):",
        sessionStart.toISOString()
      );

      const fiveMinutesInMs = 5 * 60 * 1000; // 5 minutos en milisegundos
      const timeDifference = sessionStart - now;

      console.log("Diferencia de tiempo (milisegundos):", timeDifference);
      console.log("Límite de 5 minutos (milisegundos):", fiveMinutesInMs);

      if (timeDifference <= fiveMinutesInMs && timeDifference > 0) {
        console.log(
          "Mostrando WarningModal: Faltan 5 minutos o menos para el inicio de la sesión."
        );
        setShowWarningModal(true); // Mostrar modal si faltan 5 minutos o menos
      } else {
        console.log(
          "No se muestra WarningModal: Faltan más de 5 minutos o la sesión ya comenzó."
        );
      }
    };

    checkTimeRemaining();
  }, [session_id, sessionDetails.dateStart]);

  // Manejar cambios en los selectores de pilotos
  const handleDriverChange = (position, value) => {
    setFormData((prev) => ({
      ...prev,
      [position]: value,
    }));
    // Verificar si todos los campos están completos
    setIsFormComplete(!!(formData.P1 && formData.P2 && formData.P3));
  };

  // Manejar envío del formulario (sin backend por ahora)
  const handleSubmit = (e) => {
    e.preventDefault();
    console.log("Formulario enviado con datos:", formData);
    // Aquí normalmente harías una solicitud al backend, pero por ahora solo logueamos
    navigate("/"); // Redirigir a HomePage después de "enviar"
  };

  // Determinar si es una sesión de carrera para mostrar el componente correcto
  const isRace = isRaceSession(
    sessionDetails.sessionName,
    sessionDetails.sessionType
  );

  // Cerrar el modal
  const handleCloseModal = () => {
    setShowWarningModal(false);
  };

  return (
    <div>
      <Header />
      <NavigationBar />
      <main className="pt-28 px-4">
        {/* SessionHeader con separación superior y nuevos props */}
        <SessionHeader
          countryName={sessionDetails.countryName}
          flagUrl={sessionDetails.flagUrl}
          sessionName={sessionDetails.sessionName}
          sessionType={sessionDetails.sessionType}
          className="mt-6"
        />

        {/* Isla completa para "Completa el Top 3", selectores y botón */}
        {!isRace && (
          <div className="mt-4 p-4 bg-white rounded-lg shadow-md">
            <Top3FormHeader sessionType={sessionDetails.sessionType} />

            {/* Formulario con selectores de pilotos (deshabilitado si showWarningModal es true) */}
            <form onSubmit={handleSubmit} disabled={showWarningModal}>
              <DriverSelect
                position="P1"
                value={formData.P1}
                onChange={(value) => handleDriverChange("P1", value)}
                disabled={showWarningModal}
              />
              <DriverSelect
                position="P2"
                value={formData.P2}
                onChange={(value) => handleDriverChange("P2", value)}
                disabled={showWarningModal}
              />
              <DriverSelect
                position="P3"
                value={formData.P3}
                onChange={(value) => handleDriverChange("P3", value)}
                disabled={showWarningModal}
              />

              {/* Botón "Enviar pronóstico" deshabilitado hasta que todos los campos estén completos */}
              <SubmitButton
                isDisabled={!isFormComplete || showWarningModal}
                onClick={handleSubmit}
                label="Enviar pronóstico"
                className="mt-4"
              />
            </form>
          </div>
        )}

        {/* Botón para volver a HomePage */}
        <button onClick={() => navigate("/")} className="mt-4">
          Volver
        </button>

        {/* Modal de advertencia */}
        <WarningModal isOpen={showWarningModal} onClose={handleCloseModal} />
      </main>
    </div>
  );
};

export default ProdeSessionPage;

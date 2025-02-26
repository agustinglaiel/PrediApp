import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import SessionHeader from "../components/pronosticos/SessionHeader";
import Top3FormHeader from "../components/pronosticos/Top3FormHeader";
import DriverSelect from "../components/pronosticos/DriverSelect";
import SubmitButton from "../components/pronosticos/SubmitButton";
import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";

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
  });
  const [isFormComplete, setIsFormComplete] = useState(false);

  // Simulamos un efecto para cargar datos (sin backend por ahora)
  useEffect(() => {
    // Aquí normalmente harías una solicitud al backend, pero por ahora usamos datos dummy
    console.log(`Cargando datos para sesión ${session_id}`);
  }, [session_id]);

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

        {/* Mensaje de tiempo restante (dummy por ahora) */}
        <div className="mt-4 text-orange-500">
          Puedes pronosticar hasta 5 minutos antes del evento.
        </div>

        {/* Isla completa para "Completa el Top 3", selectores y botón */}
        {!isRace && (
          <div className="mt-4 p-4 bg-white rounded-lg shadow-md">
            <Top3FormHeader sessionType={sessionDetails.sessionType} />

            {/* Formulario con selectores de pilotos */}
            <form onSubmit={handleSubmit}>
              <DriverSelect
                position="P1"
                value={formData.P1}
                onChange={(value) => handleDriverChange("P1", value)}
              />
              <DriverSelect
                position="P2"
                value={formData.P2}
                onChange={(value) => handleDriverChange("P2", value)}
              />
              <DriverSelect
                position="P3"
                value={formData.P3}
                onChange={(value) => handleDriverChange("P3", value)}
              />

              {/* Botón "Enviar pronóstico" deshabilitado hasta que todos los campos estén completos */}
              <SubmitButton
                isDisabled={!isFormComplete}
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
      </main>
    </div>
  );
};

export default ProdeSessionPage;

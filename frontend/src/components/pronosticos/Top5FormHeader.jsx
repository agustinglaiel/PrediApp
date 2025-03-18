import React from "react";

const Top5FormHeader = ({ sessionType }) => {
  return (
    <div className="p-4 bg-white rounded-lg mb-4">
      <h2 className="text-lg font-bold">Completa el Top 5</h2>
      <p className="text-sm text-gray-600">
        Cada pleno te otorga 3 puntos. Si se acierta dentro del top, se otorga 1
        punto.
      </p>
    </div>
  );
};

export default Top5FormHeader;

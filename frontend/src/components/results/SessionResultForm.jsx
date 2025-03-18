// frontendnuevo/src/components/results/SessionResultForm.jsx
import React from "react";
import Top3FormHeader from "../pronosticos/Top3FormHeader";

const SessionResultForm = ({ sessionType, drivers }) => {
  return (
    <div className="mt-4 p-4 bg-white rounded-lg shadow-md">
      <Top3FormHeader sessionType={sessionType} />
      <div className="flex flex-col gap-4 mt-4">
        {/* P1 */}
        <div className="mb-4 ml-4">
          <label className="block text-sm font-medium text-black">P1</label>
          <div className="mt-1 block w-full py-2 px-3 text-gray-700 bg-gray-100 border border-gray-300 rounded-md">
            {drivers.p1?.full_name || "No seleccionado"}
          </div>
        </div>
        {/* P2 */}
        <div className="mb-4 ml-4">
          <label className="block text-sm font-medium text-black">P2</label>
          <div className="mt-1 block w-full py-2 px-3 text-gray-700 bg-gray-100 border border-gray-300 rounded-md">
            {drivers.p2?.full_name || "No seleccionado"}
          </div>
        </div>
        {/* P3 */}
        <div className="mb-4 ml-4">
          <label className="block text-sm font-medium text-black">P3</label>
          <div className="mt-1 block w-full py-2 px-3 text-gray-700 bg-gray-100 border border-gray-300 rounded-md">
            {drivers.p3?.full_name || "No seleccionado"}
          </div>
        </div>
      </div>
    </div>
  );
};

export default SessionResultForm;

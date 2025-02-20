import axios from "axios";

const API_URL = "http://localhost:8055";

export const getResults = async (sessionID) => {
  try {
    const response = await axios.get(`${API_URL}/results/session/${sessionID}`);
    return response.data;
  } catch (error) {
    throw new Error(error.response.data.message || "Error fetching results.");
  }
};

import { render, screen } from "@testing-library/react";
import App from "./App";

test("renders welcome message", () => {
  render(<App />);
  const linkElement = screen.getByText(/Bienvenido/i); // Asegúrate de que este texto existe en tu App
  expect(linkElement).toBeInTheDocument();
});

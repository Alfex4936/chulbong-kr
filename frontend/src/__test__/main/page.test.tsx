import "@testing-library/jest-dom";
import { render, screen, waitFor } from "@testing-library/react";
import Page from "@/app/(non-navbar)/page";

HTMLCanvasElement.prototype.getContext = jest.fn();
describe("main page", () => {
  it("should render the button with the text 철봉 지도 바로 가기", async () => {
    render(<Page />);

    await waitFor(() => {});

    const navBtn = screen.getByTestId("nav-btn");

    expect(navBtn).toBeInTheDocument();
    expect(navBtn).toHaveTextContent("철봉 지도 바로 가기");
  });
});

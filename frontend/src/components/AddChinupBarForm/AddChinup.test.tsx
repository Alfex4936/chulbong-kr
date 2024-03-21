import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import "@testing-library/jest-dom";
import { render, screen } from "@testing-library/react";
import AddChinupBarForm from "./AddChinupBarForm";

const mockMarker = {
  setPosition: jest.fn(),
  getPosition: jest.fn().mockReturnValue({
    La: 37.5665,
    Ma: 126.978,
    getLat: () => 37.5665,
    getLng: () => 126.978,
  }),
  setImage: jest.fn(),
  setMap: jest.fn(),
  Gb: "someId",
};

const mockMap = {
  getCenter: jest.fn(),
  setLevel: jest.fn(),
  setCenter: jest.fn().mockReturnValue({
    getLat: () => 37.5665,
    getLng: () => 126.978,
  }),
  getLevel: jest.fn().mockReturnValue(3),
};

const mockProps = {
  setState: jest.fn(),
  setIsMarked: jest.fn(),
  setMarkerInfoModal: jest.fn(),
  setCurrentMarkerInfo: jest.fn(),
  setMarkers: jest.fn(),
  map: mockMap,
  marker: mockMarker,
  markers: [],
  clusterer: {
    addMarker: jest.fn(),
    removeMarker: jest.fn(),
    addMarkers: jest.fn(),
    clear: jest.fn(),
    redraw: jest.fn(),
  },
};

jest.mock("nanoid", () => {
  return {
    nanoid: () => {},
  };
});

const queryClient = new QueryClient();

describe("AddChinupBarForm Component", () => {
  test("renders without crashing", () => {
    render(
      <QueryClientProvider client={queryClient}>
        <AddChinupBarForm {...mockProps} />
      </QueryClientProvider>
    );

    const titleText = screen.getByText("위치 등록");
    expect(titleText).toBeInTheDocument();
  });
});

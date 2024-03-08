import styled from "@emotion/styled";
import { keyframes } from "@emotion/react";

const shimmer = keyframes`
  0% {
    background-position: -468px 0;
  }
  100% {
    background-position: 468px 0;
  }
`;

export const Container = styled.div`
  max-height: 200px;

  overflow: auto;

  & > div:last-of-type {
    border-bottom: none;

    margin-bottom: 0.5rem;
  }
`;

export const ListContainer = styled.div``;

export const RangeContainer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;

  padding: 1rem;
`;

export const SearchButtonContainer = styled.div`
  width: 30px;
  padding: 0 1rem;

  & > button {
    font-size: 0.7rem;
    margin: 0;
  }
`;

export const LoadList = styled.div``;

export const MarkerList = styled.div`
  display: flex;
  align-items: center;

  padding: 1rem;
  margin: 0 auto 1rem auto;

  width: 250px;

  border-radius: 0.4rem;
  background-color: #e9efff;
`;

export const ListSkeleton = styled.div`
  display: flex;
  align-items: center;

  margin: 0 auto 1rem auto;

  padding: 1rem;

  height: 57px;
  width: 250px;

  background: #f6f7f8;
  background-image: linear-gradient(
    90deg,
    #f0f0f0 25%,
    #f7f7f7 50%,
    #f0f0f0 75%
  );

  animation: ${shimmer} 1.2s ease-in-out infinite;
  border-radius: 0.4rem;
`;

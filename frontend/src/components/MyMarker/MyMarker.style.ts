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

  padding: 1rem;

  & > div:last-of-type {
    border-bottom: none;

    margin-bottom: 0.5rem;
  }
`;

export const ListContainer = styled.div``;

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

export const MarkerListTop = styled.div`
  flex-grow: 1;
`;

export const AddressText = styled.p`
  font-size: 0.7rem;
  color: #777;

  text-align: left;
`;

export const ListSkeleton = styled.div`
  display: flex;
  align-items: center;

  margin: 1rem auto 1rem auto;

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

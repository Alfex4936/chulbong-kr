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
  // border: 1px solid red;

  height: 200px;

  overflow: auto;
`;

export const LoadList = styled.div`
  // border: 1px solid red;

  // height: 30px;
`;

export const MarkerList = styled.div`
  display: flex;
  align-items: center;

  padding: 1rem;
  margin: auto;

  width: 250px;

  border-bottom: 1px solid #ccc;
`;

export const ListSkeleton = styled.div`
  display: flex;
  align-items: center;

  padding: 1rem;

  height: 60px;
  width: 250px;

  border-bottom: 1px solid #ccc;

  & > div:first-of-type {
    height: 24px;
    width: 130px;

    background: #f6f7f8;
    background-image: linear-gradient(
      90deg,
      #f0f0f0 25%,
      #f7f7f7 50%,
      #f0f0f0 75%
    );

    animation: ${shimmer} 1.2s ease-in-out infinite;

    border-radius: 1rem;
  }

  & > div:last-of-type {
    height: 22px;
    width: 22px;

    background: #f6f7f8;
    background-image: linear-gradient(
      90deg,
      #f0f0f0 25%,
      #f7f7f7 50%,
      #f0f0f0 75%
    );

    animation: ${shimmer} 1.2s ease-in-out infinite;

    border-radius: 50%;
  }
`;

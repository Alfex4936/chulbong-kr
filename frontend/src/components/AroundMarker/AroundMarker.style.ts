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
  max-height: 400px;

  overflow: auto;

  scroll-behavior: smooth;

  z-index: 100;

  & > div:last-of-type {
    border-bottom: none;

    margin-bottom: 0.5rem;
  }
`;

export const MessageRed = styled.p`
  font-size: 0.8rem;
  color: #ff6060;

  padding: 0 1rem;
`;

export const ListContainer = styled.div``;

export const RangeContainer = styled.div`
  position: sticky;
  top: 0;
  left: 0;

  display: flex;
  justify-content: space-between;
  align-items: center;

  padding: 1rem;
  margin: auto;

  background-color: white;

  width: 90%;

  z-index: 1;

  & > p:first-of-type {
    font-size: 0.8rem;
    font-weight: bold;
  }

  & div > input {
    width: 90%;

    border: 1px solid red;
  }
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

  width: 90%;

  border-radius: 0.4rem;
  background-color: #e9efff;
`;

export const MarkerListTop = styled.div`
  flex-grow: 1;
`;

export const DescriptionWrap = styled.div`
  display: flex;
  align-items: center;
`;

export const AddressText = styled.p`
  font-size: 0.7rem;
  color: #777;

  text-align: left;
`;

export const Distance = styled.div`
  font-size: 0.7rem;
  margin-right: 0.5rem;
  font-weight: bold;
`;

export const Description = styled.p`
  max-width: 150px;
  text-align: left;

  white-space: nowrap;

  overflow: hidden;

  text-overflow: ellipsis;
`;

export const ListSkeleton = styled.div`
  display: flex;
  align-items: center;

  margin: 0 auto 1rem auto;

  padding: 1rem;

  height: 57px;
  width: 90%;

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

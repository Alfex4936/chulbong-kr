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
  display: flex;
  flex-direction: column;

  // margin-top: -1rem;
`;

export const Weather = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;

  margin-bottom: 0.5rem;

  & > div:first-of-type {
    font-size: 1.3rem;
    font-weight: bold;

    margin-left: 0.5rem;
  }
  // & > div:last-of-type {
  //   font-size: 0.8rem;
  // }
`;

export const WeatherSkeleton = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;

  margin: auto;
  margin-bottom: 0.5rem;

  width: 60px;
  height: 40px;

  background: #f6f7f8;
  background-image: linear-gradient(
    90deg,
    #f0f0f0 25%,
    #f7f7f7 50%,
    #f0f0f0 75%
  );

  animation: ${shimmer} 1.2s ease-in-out infinite;

  border-radius: 1rem;
`;

export const ImagesContainer = styled.div`
  display: flex;
  justify-content: flex-start;
  margin-bottom: 1.5rem;
`;

export const imageWrap = styled.div`
  position: relative;

  width: 85%;

  & > div:first-of-type {
    max-height: 360px;

    border-radius: 1rem;
    overflow: hidden;
  }

  & img {
    object-fit: cover;
    background-position: center;
    background-size: cover;

    width: 100%;
    max-height: 360px;
    display: block;

    box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

    user-select: none;
  }

  margin: auto;
`;

export const ImagePreviewWrap = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;

  width: 15%;

  & > button {
    background: transparent;
    padding: 0;
    cursor: pointer;

    border: none;
    border-radius: 0.4rem;
    margin: 0.3rem;
    box-shadow: rgba(0, 0, 0, 0.24) 0px 3px 8px;

    width: 80%;

    overflow: hidden;

    & > img {
      width: 100%;
      display: block;
    }
  }
`;

export const description = styled.div`
  position: absolute;

  bottom: 0;
  left: 0;

  font-size: 1.3rem;
  margin-top: 1rem;
  color: #fff;

  width: 100%;

  background-color: rgba(0, 0, 0, 0.5);

  padding: 1rem;

  border-bottom-left-radius: 1rem;
  border-bottom-right-radius: 1rem;

  & > div {
    white-space: nowrap;
    overflow: hidden;

    text-overflow: ellipsis;
  }

  &:hover > div {
    max-height: 200px;
    word-wrap: break-word;
    white-space: -moz-pre-wrap;
    white-space: pre-wrap;

    overflow: auto;

    text-overflow: none;
  }
`;

export const AddressText = styled.div`
  margin-bottom: 0.5rem;

  font-weight: bold;

  & > div:first-of-type {
    font-size: 0.8rem;
    font-weight: 400;

    color: #777;
  }
`;

export const BottomButtons = styled.div`
  margin-bottom: -1rem;
`;

export const DislikeCount = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;

  position: absolute;
  top: -7px;
  right: -10px;

  width: 20px;
  height: 13px;

  font-size: 0.5rem;
  color: #fff;

  border-radius: 10px;

  background-color: #ff7e7e;
`;

export const InputWrap = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
`;

export const ButtonWrap = styled.div`
  width: 100%;
  & button {
    margin: 0;

    width: 30px;

    margin: 0.5rem 0.5rem 0 0.5rem;
  }
`;

export const DescInput = styled.input`
  background-color: rgba(0, 0, 0, 0.3);
  color: #fff;

  outline: none;
  border: 1px solid #ccc;
  border-radius: 4px;

  width: 70%;
  height: 25px;

  font: inherit;

  font-size: 0.8rem;
`;

export const ErrorBox = styled.div`
  text-align: center;

  font-size: 0.7rem;

  color: red;
`;

export const Facilities = styled.div`
  display: flex;
  justify-content: center;

  margin-bottom: 0.5rem;

  & > div {
    margin: 0 0.5rem;
    display: flex;
    align-items: center;

    & > span:first-of-type {
      font-size: 1.2rem;
      margin-right: 0.3rem;
    }
    & > span:last-of-type {
      font-size: 0.9rem;
    }
  }
`;

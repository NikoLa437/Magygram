import React, { createContext, useReducer } from "react";
import { storyReducer } from "../reducers/StoryReducer";

export const StoryContext = createContext();

const StoryContextProvider = (props) => {
	const [storyState, dispatch] = useReducer(storyReducer, {
		createStory: {
			showModal: false,
			showError: false,
			errorMessage: "",
			showSuccessMessage: false,
			successMessage: "",
		},
		storyline: {
			stories: [],
		},
		highlights: {
			showModal: false,
			showError: false,
			errorMessage: "",
			showHighlightsName: false,
			stories: [],
		},
		highlightsSliderModal: {
			showModal: false,
			highlights: [],
		},
		profileHighlights: {
			highlights: [],
		},
		storySliderModal: {
			showModal: false,
			stories: [],
			firstUnvisitedStory: 0,
		},
		iHaveAStory: true,
	});

	return <StoryContext.Provider value={{ storyState, dispatch }}>{props.children}</StoryContext.Provider>;
};

export default StoryContextProvider;

import React, { createContext, useReducer } from "react";
import { notificationReducer } from "../reducers/NotificationReducer";

export const NotificationContext = createContext();

const NotificationContextProvider = (props) => {
	const [notificationState, dispatch] = useReducer(notificationReducer, {
		userFollowRequests: {
			userInfos: [],
		},
		notifications: [],
		notificationsNumber: 0,
	});

	return <NotificationContext.Provider value={{ notificationState, dispatch }}>{props.children}</NotificationContext.Provider>;
};

export default NotificationContextProvider;

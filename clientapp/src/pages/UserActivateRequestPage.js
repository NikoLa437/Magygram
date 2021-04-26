import React, { useContext, useEffect } from "react";
import { UserContext } from "../contexts/UserContext";
import { userService } from "../services/UserService";


const UserActivateRequestPage = (props) => {

    const {userState, dispatch } = useContext(UserContext);

    const userId = props.match.params.id;
    
    const handleSubmit = (e) => {
        e.preventDefault();
        let resendActivationLink = {
			email: userState.blockedUser.emailAddress
		};

        userService.resendActivationLink(resendActivationLink,dispatch)
    };

    useEffect(() => {
        userService.checkIfUserIdExist(userId,dispatch)
    },[userState.blockedUser.emailAddress])

	return (
            <React.Fragment>
            <section className="login-clean">
                <div className="illustration">
                    <i className="icon ion-ios-navigate"></i>
                </div>
                <form method="post" onSubmit={handleSubmit}>
                    <div>Vas nalog je privremeno blokiran. Ukoliko zelite ponovno aktiviranje naloga pritisnite na dugme ispod nakon cega ce Vam na email: <b>{userState.blockedUser.emailAddress}</b> stici aktivacioni link.</div>

                    <div className="form-group">
                        <input className="btn btn-primary btn-block" type="submit" value="Send activation mail" />
                    </div>
                </form>
			</section>
		</React.Fragment>
	);
};

export default UserActivateRequestPage;

import React, { useState, useContext, useEffect } from "react";
import { userService } from "../services/UserService";
import Axios from "axios";
import { authHeader } from "../helpers/auth-header";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";

const EditProfile = () => {
	// history = useHistory();
	const navStyle = { height: "50px", borderBottom: "1px solid rgb(200,200,200)" };
	const inputStyle = { border: "1px solid rgb(200,200,200)", color: "rgb(210,210,210)", textAlign: "center" };
	const iconStyle = { fontSize: "30px", margin: "0px", marginLeft: "13px" };
	//const iconStyle1 = { fontSize: "30px", margin: "0px", marginLeft: "200px" };
	const imgStyle = { left: "0", width: "30px", height: "30px", marginLeft: "13px", borderWidth: "1px", borderStyle: "solid" };
	// const imgProfileStyle = { left: "20", width: "150px", height: "150px", marginLeft: "100px", borderWidth: "1px", borderStyle: "solid" };
	// const nameStyle = { left: "20",  marginLeft: "13px"}
	// const editStyle = {color: "black", left: "20",  marginLeft: "13px",marginRight: "13px", borderWidth: "1px", borderStyle: "solid" }
	const sectionStyle = { left: "20", marginLeft: "100px" };

	const [email, setEmail] = useState("");
	const [id, setId] = useState("");
	const [name, setName] = useState("");
	const [username, setUsername] = useState("");
	const [bio, setBio] = useState("");
	const [img, setImg] = useState("");
	const [website, setWebsite] = useState("");
	const [number, setNumber] = useState("");

	const [emailInput, setEmailInput] = useState("");
	const [nameInput, setNameInput] = useState("");
	const [usernameInput, setUsernameInput] = useState("");
	const [bioInput, setBioInput] = useState("");
	//const [imgInput, setImgInput] = useState("");
	const [websiteInput, setWebsiteInput] = useState("");
	const [numberInput, setNumberInput] = useState("");

	useEffect(() => {
		Axios.get(`https://localhost:460/api/users/logged`, { validateStatus: () => true, headers: authHeader() })
			.then((res) => {
				setId(res.data.id);
				console.log(res.data);
				if (res.data.imageUrl == "") setImg("assets/img/profile.jpg");
				else setImg(res.data.imageUrl);

				Axios.get(`https://localhost:460/api/users/` + res.data.id, { validateStatus: () => true, headers: authHeader() })
					.then((res) => {
						console.log(res.data);
						setName(res.data.Name);
						setUsername(res.data.Username);
						setBio(res.data.Bio);
						setEmail(res.data.Email);
						setWebsite(res.data.Website);
						setNumber(res.data.Number);
					})
					.catch((err) => {
						console.log(err);
					});
			})
			.catch((err) => {
				console.log(err);
			});
	});

	const handleSettings = () => {
		alert("TOD1O");
	};

	const handleSubmit = (e) => {
		e.preventDefault();

		let user = {
			id,
			nameInput,
			usernameInput,
			emailInput,
			bioInput,
			websiteInput,
			numberInput,
		};
	};

	const handleLogout = () => {
		userService.logout();
	};

	const handleProfile = () => {
		window.location = `#/profile`;
	};

	const handleChange = (e) => {
		setName(e.target.value);
	};

	return (
		<React.Fragment>
			<div>
				<nav className="navbar navbar-light navbar-expand-md navigation-clean" style={navStyle}>
					<div className="container">
						<div>
							<img src="assets/img/logotest.png" alt="NistagramLogo" />
						</div>
						<button className="navbar-toggler" data-toggle="collapse">
							<span className="sr-only">Toggle navigation</span>
							<span className="navbar-toggler-icon"></span>
						</button>
						<div>
							<input type="text" style={inputStyle} placeholder="Search" value="Search" />
						</div>
						<div className="d-xl-flex align-items-xl-center dropdown">
							<i className="fa fa-home" style={iconStyle} />
							<i className="la la-wechat" style={iconStyle} />
							<i className="la la-compass" style={iconStyle} />
							<i className="fa fa-heart-o" style={iconStyle} />
							<img className="rounded-circle dropdown-toggle" data-toggle="dropdown" style={imgStyle} src={img} alt="ProfileImage" />
							<ul style={{ width: "200px", marginLeft: "15px" }} class="dropdown-menu">
								<li>
									<button className="la la-user btn shadow-none" onClick={handleProfile}>
										{" "}
										Profile
									</button>
								</li>
								<li>
									<button className="la la-cog btn shadow-none" onClick={handleSettings}>
										{" "}
										Settings
									</button>
								</li>
								<hr className="solid" />
								<li>
									<button className=" btn shadow-none" onClick={handleLogout}>
										{" "}
										Logout
									</button>
								</li>
							</ul>
						</div>
						<div>{name}</div>
					</div>
				</nav>
				<br />
				<br />
				<div className="container">
					<form method="post" onSubmit={handleSubmit}>
						<div>
							<h2 className="text-center">
								<strong>Edit</strong> profile
							</h2>
							<br />
							<div className="form-group">
								<text>Name</text>
								<input className="form-control" required type="text" name="name" placeholder={name} value={nameInput} onChange={(e) => setNameInput(e.target.value)} />
							</div>

							<div className="form-group">
								<text>Email</text>
								<input className="form-control" required type="email" name="email" placeholder={email} value={emailInput} onChange={(e) => setEmailInput(e.target.value)} />
							</div>

							<div className="form-group">
								<text>Username</text>
								<input
									className="form-control"
									required
									type="username"
									name="username"
									placeholder={username}
									value={usernameInput}
									onChange={(e) => setUsernameInput(e.target.value)}
								/>
							</div>

							<div className="form-group">
								<text>Website</text>
								<input className="form-control" required type="text" name="websiteInput" placeholder={website} value={websiteInput} onChange={(e) => setWebsiteInput(e.target.value)} />
							</div>

							<div className="form-group">
								<text>Bio</text>
								<input className="form-control" required type="text" name="bioInput" placeholder={bio} value={bioInput} onChange={(e) => setBioInput(e.target.value)} />
							</div>

							<div className="form-group">
								<text>Number</text>
								<input className="form-control" required type="text" name="numberInput" placeholder={number} value={numberInput} onChange={(e) => setNumberInput(e.target.value)} />
							</div>

							<div class="flexbox-container">
								<div>Date of birth</div>
								<div style={sectionStyle}>
									<DatePicker />
								</div>
							</div>
							<br />
							<div class="flexbox-container">
								<div>Gender</div>
								<div style={sectionStyle}>
									<select id="dropdown">
										<option value="1"> Male</option>
										<option value="2"> Female</option>
									</select>
								</div>
							</div>
							<br />
							<div className="form-group">
								<input className="btn btn-primary btn-block" type="submit" value="Save" />
							</div>
						</div>
					</form>
				</div>
			</div>
		</React.Fragment>
	);
};

export default EditProfile;

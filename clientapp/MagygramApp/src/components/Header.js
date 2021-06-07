import React, { useEffect, useState, useContext } from "react";
import { userService } from "../services/UserService";
import Axios from "axios";
import { Link, useHistory } from "react-router-dom";
import { authHeader } from "../helpers/auth-header";
import AsyncSelect from "react-select/async";
import { searchService } from "../services/SearchService";
import { postService } from "../services/PostService";
import { PostContext } from "../contexts/PostContext";

const Header = () => {
	const { dispatch } = useContext(PostContext);

	const history = useHistory();
	const navStyle = { height: "50px", borderBottom: "1px solid rgb(200,200,200)" };
	const iconStyle = { fontSize: "30px", margin: "0px", marginLeft: "13px" };
	const imgStyle = { width: "30px", height: "30px", marginLeft: "13px", borderWidth: "1px", borderStyle: "solid" };

	const [name, setName] = useState("");
	const [img, setImg] = useState("");
	const [currentId, setCurrentId] = useState();
	const [search, setSearch] = useState("");

	useEffect(() => {
		Axios.get(`/api/users/logged`, { validateStatus: () => true, headers: authHeader() })
			.then((res) => {
				setCurrentId(res.data.id);
				if (res.data.imageUrl == "") setImg("assets/img/profile.jpg");
				else setImg(res.data.imageUrl);
			})
			.catch((err) => {
				console.log(err);
			});
	});

	const loadOptions = (value, callback) => {
		if(value.startsWith('#')){
            setTimeout(() => {
                searchService.guestSearchHashtagPosts(value,callback)
            }, 1000);
        }else if(value.startsWith('%')){
			setTimeout(() => {
                searchService.guestSearchLocation(value,callback)
            }, 1000);
		}else{
            setTimeout(() => {
                searchService.userSearchUsers(value,callback)
            }, 1000);
        }
	  };

	const onInputChange = (inputValue, { action }) => {
		switch (action) {
			case "set-value":
				return;
			case "menu-close":
				setSearch("");
				return;
			case "input-change":
				setSearch(inputValue);
				return;
			default:
				return;
		}
	};

	const onChange = (option) => {
		if (option.searchType === "hashtag") {
			postService.findPostsForUserByHashtag(option.value,dispatch);
		}else if(option.searchType === "location"){
			postService.findPostsForUserByLocation(option.value,dispatch);
		} else {
			window.location = "#/user/" + option.id;
		}

		return false;
	};

	const handleLogout = () => {
		userService.logout();
	};

	const handleProfile = () => {
		let path = `/profile`;
		history.push(path);
	};

	const handleSettings = () => {
		alert("TOD1O");
	};

	return (
		<nav className="navbar navbar-light navbar-expand-md navigation-clean" style={navStyle}>
			<div className="container">
				<div>
					<img src="assets/img/logotest.png" alt="NistagramLogo" />
				</div>
				<button className="navbar-toggler" data-toggle="collapse">
					<span className="sr-only">Toggle navigation</span>
					<span className="navbar-toggler-icon"></span>
				</button>
				<div style={{ width: "300px" }}>
					<AsyncSelect defaultOptions loadOptions={loadOptions} onInputChange={onInputChange} onChange={onChange} placeholder="search" inputValue={search} />
				</div>
				<div className="d-xl-flex align-items-xl-center dropdown">
					<i className="fa fa-home" style={iconStyle} />
					<i className="la la-wechat" style={iconStyle} />
					<i className="la la-compass" style={iconStyle} />
					<i className="fa fa-heart-o" style={iconStyle} />
					<img className="rounded-circle dropdown-toggle" data-toggle="dropdown" style={imgStyle} src={img} alt="ProfileImage" />
					<ul style={{ width: "200px", marginLeft: "15px" }} class="dropdown-menu">
						<li>
							<Link className="la la-user btn shadow-none" to={"/profile?userId=" + localStorage.getItem("userId")}>
								Profile
							</Link>
						</li>
						<li>
							<button className="la la-cog btn shadow-none" onClick={handleSettings}>
								Settings
							</button>
						</li>
						<hr className="solid" />
						<li>
							<button className=" btn shadow-none" onClick={handleLogout}>
								Logout
							</button>
						</li>
					</ul>
				</div>
				<div>{name}</div>
			</div>
		</nav>
	);
};

export default Header;

import React, {Fragment, useEffect, useState} from "react";
import {BrowserRouter, Link, Route, Routes} from "react-router-dom";
import {useCookie} from "react-use";

import Home from "./Home";
import Movie from "./Movie";
import MovieAddEdit from "./MovieAddEdit";
import Movies from "./Movies";
import Genres from "./Genres"
import Genre from "./Genre"
import Login from "./Login"
import Logout from "./Logout";
import Admin from "./Admin";

export default function App() {
    let [value, updateCookie, deleteCookie] = useCookie("Token")
    const [forbidden, setForbidden] = useState(true)

    let checkStatus = () => {
        const headers = new Headers()
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", "Bearer " + value)

        const requestOptions = {
            method: "POST",
            headers: headers
        }

        fetch(`${process.env.REACT_APP_API_URL}/${process.env.REACT_APP_API_VERSION}/checkup`, requestOptions)
            .then((response) => {
                if (response.status !== 200)
                    return true
                else
                    return false
            })
    }

    useEffect(() => {
        const headers = new Headers()
        headers.append("Content-Type", "application/json")
        headers.append("Authorization", "Bearer " + value)

        const requestOptions = {
            method: "POST",
            headers: headers
        }

        fetch(`${process.env.REACT_APP_API_URL}/${process.env.REACT_APP_API_VERSION}/checkup`, requestOptions)
            .then((response) => {
                if (response.status !== 200)
                    setForbidden(true)
                else
                    setForbidden(false)
            })
    },  [value])

    return (
        <BrowserRouter>
            <nav className="navbar navbar-expand sticky-top navbar-light"
                 style={{backgroundColor: `rgba(37, 167, 239, 0.62)`}}>
                <div className="container-fluid">
                    <Link className={"navbar-brand"} to={"/"}>React/Go Movies</Link>
                    <div className={"collapse navbar-collapse"}>
                        <ul className={"navbar-nav"}>
                            <li className={"nav-item"}>
                                <Link className={"btn btn-outline-dark mx-2"} to={"/movies"}>Movies</Link>
                            </li>
                            <span className={"divider"}></span>
                            <li className={"nav-item"}>
                                <Link className={"btn btn-outline-dark mx-2"} to={"/genres"}>Genres</Link>
                            </li>
                            {!forbidden && (
                                <Fragment>
                                    <li className="nav-item">
                                        <Link className={"btn btn-outline-dark mx-2"} to={"/admin/movie/0"}>Add a
                                            Movie</Link>
                                    </li>
                                    <li className="nav-item">
                                        <Link className={"btn btn-outline-dark mx-2"} to={"/admin"}>Manage Movies</Link>
                                    </li>
                                </Fragment>
                            )}
                        </ul>
                    </div>
                    <div>
                        {value === null && (
                            <Link className={"btn btn-outline-dark"} to={"/login"}>Login</Link>
                        )}
                        {value !== null && (
                            <Link className={"btn btn-outline-dark"} to={"/logout"}>Logout</Link>
                        )}
                    </div>
                </div>
            </nav>
            <br/>
            <div className={"container"}>
                <div className={"row"}>
                    <div>
                        <Routes>
                            <Route path={"/"} element={<Home/>}/>
                            <Route path={"/movies"} element={<Movies/>}/>
                            <Route path={"/movies/:id"} element={<Movie/>}/>
                            <Route path={"/genres"} element={<Genres/>}/>
                            <Route path={"/genres/:id"} element={<Genre/>}/>
                            <Route path={"/login"}
                                   element={<Login updateCookie={updateCookie}/>}/>
                            <Route path={"/logout"} element={<Logout deleteCookie={deleteCookie}/>}/>
                            <Route path={"/admin"} element={<Admin jwt={value} checkStatus={checkStatus}/>}/>
                            <Route path={"/admin/movie/:id"}
                                   element={<MovieAddEdit jwt={value} checkStatus={checkStatus}/>}/>
                        </Routes>
                    </div>
                </div>
            </div>
        </BrowserRouter>
    )
}
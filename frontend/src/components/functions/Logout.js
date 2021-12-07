import React, {useEffect} from "react";
import {useNavigate} from "react-router-dom";

export default function Logout({deleteCookie}) {
    const navigate = useNavigate()

    useEffect(() => {
        deleteCookie()
        navigate("/")
    }, [deleteCookie, navigate])

    return (<h2>Logout</h2>)
}
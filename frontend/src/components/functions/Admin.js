import React, {Fragment, useEffect, useState} from "react";
import {Link, useNavigate} from "react-router-dom";

export default function Admin({jwt, forbidden}) {
    const navigate = useNavigate()

    const [movies, setMovies] = useState([])
    const [error, setError] = useState(null)

    useEffect(() => {
        if (jwt === null || forbidden) {
            navigate('/login')
            return
        }

        fetch(`${process.env.REACT_APP_API_URL}/${process.env.REACT_APP_API_VERSION}/movies`)
            .then((response) => {
                if (response.status !== 200)
                    setError("Invalid response code: ", response.status)
                else
                    setError(null)
                return response.json()
            })
            .then((json) => {
                setMovies(json.movies)
            })
    }, [forbidden, jwt, navigate])

    if (error !== null)
        return <div>Error: {error.message}</div>
    else {
        return (
            <Fragment>
                <h2>Choose a Movie</h2>
                <div className={"list-group"}>
                    {movies.map((m, index) => (
                        <Link to={`/admin/movie/${m.id}`}
                              className={"list-group-item list-group-item-action"} key={index}>{m.title}</Link>
                    ))}
                </div>
            </Fragment>
        )
    }
}
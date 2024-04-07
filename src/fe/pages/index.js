import Head from 'next/head';
import styles from '../styles/Home.module.css';
import {useEffect, useState} from "react";

export default function Home() {
    const [serverStatus, setServerStatus] = useState('loading');
    const [result, setResult] = useState('');
    const [showResult, setShowResult] = useState(false);

    const [startNode, setStartNode] = useState('');
    const [endNode, setEndNode] = useState('');


    useEffect(() => {
        fetchServerStatus().then(r => console.log('Server status fetched'));
    }, []);

    // Fetch the server status
    // Use async await to fetch the data
    const fetchServerStatus = async () => {
        try {
            const response = await fetch('http://localhost:8000/status');
            if (!response.ok) { // if HTTP-status is 200-299
                // get the error message from the server,
                const error = await response.text();
                throw new Error(error);
            }
            const data = await response.json();
            setServerStatus(data.status);
        } catch (error) {
            console.error('Failed to fetch server status: ', error);
        }
    };

    const testBFS = async () => {
        // Create a new URL object
        const url = new URL('http://localhost:8000/bfs');

        // Add the startNode and endNode as query parameters
        url.searchParams.append('start', startNode);
        url.searchParams.append('end', endNode);

        try {
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
            });

            if (!response.ok) { // if HTTP-status is 200-299
                // get the error message from the server,
                const error = await response.text();
                throw new Error(error);
            }

            const data = await response.json();
            console.log('Success:', data);
            setResult(data.result);
            setShowResult(true);
        } catch (error) {
            console.error('Error:', error);
        }
    }

    // Call the fetchServerStatus function


    function testDFS() {
        // Create a new URL object
        const url = new URL('http://localhost:8000/dfs');

        // Add the startNode and endNode as query parameters
        url.searchParams.append('start', startNode);
        url.searchParams.append('end', endNode);

        fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        })
            .then(response => response.json())
            .then(data => {
                console.log('Success:', data);
                setResult(data.result);
                setShowResult(true);
            })
            .catch((error) => {
                console.error('Error:', error);
            });
    }

    return (
        <div className={styles.container}>
          <Head>
            <title>Create Next App</title>
            <link rel="icon" href="/favicon.ico" />
          </Head>

          <main>
            <h1 className={styles.title}>
              Hello Next js & GO!
            </h1>

            <p className={styles.description}>
              Server status : {serverStatus}
            </p>

              <div className={styles.grid}>

                  <div className={styles.card}>
                      <input type="text" id="startNode" placeholder="Start Node"
                             onChange={(e) => setStartNode(e.target.value)}></input>
                      <input type="text" id="endNode" placeholder="End Node"
                             onChange={(e) => setEndNode(e.target.value)}></input>
                      <button onClick={testBFS}>Test BFS</button>
                      <button onClick={testDFS}>Test DFS</button>
                  </div>

                  {showResult ? result : null}

              </div>
          </main>

            <footer>
                <a
                    href="https://vercel.com?utm_source=create-next-app&utm_medium=default-template&utm_campaign=create-next-app"
                    target="_blank"
                    rel="noopener noreferrer"
                >
                    Powered by{' '}
                    <img src="/vercel.svg" alt="Vercel" className={styles.logo} />
            </a>
          </footer>

          <style jsx>{`
            main {
              padding: 5rem 0;
              flex: 1;
              display: flex;
              flex-direction: column;
              justify-content: center;
              align-items: center;
            }
            footer {
              width: 100%;
              height: 100px;
              border-top: 1px solid #eaeaea;
              display: flex;
              justify-content: center;
              align-items: center;
            }
            footer img {
              margin-left: 0.5rem;
            }
            footer a {
              display: flex;
              justify-content: center;
              align-items: center;
              text-decoration: none;
              color: inherit;
            }
            code {
              background: #fafafa;
              border-radius: 5px;
              padding: 0.75rem;
              font-size: 1.1rem;
              font-family:
                Menlo,
                Monaco,
                Lucida Console,
                Liberation Mono,
                DejaVu Sans Mono,
                Bitstream Vera Sans Mono,
                Courier New,
                monospace;
            }
          `}</style>

          <style jsx global>{`
            html,
            body {
              padding: 0;
              margin: 0;
              font-family:
                -apple-system,
                BlinkMacSystemFont,
                Segoe UI,
                Roboto,
                Oxygen,
                Ubuntu,
                Cantarell,
                Fira Sans,
                Droid Sans,
                Helvetica Neue,
                sans-serif;
            }
            * {
              box-sizing: border-box;
            }
          `}</style>
        </div>
    );
}

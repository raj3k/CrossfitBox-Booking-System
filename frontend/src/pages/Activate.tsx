import { useParams } from 'react-router-dom';

const Activate: React.FC = () => {
    const { uuid } = useParams();

    return (
        <h1>Activate User { uuid } Page</h1>
    )
}

export default Activate;
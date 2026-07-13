import React from 'react';
import { Route, Routes } from 'react-router-dom';
import { AppRootProps } from '@grafana/data';
import { ROUTES } from '../../constants';
import AnalysisPage from '../../pages/AnalysisPage';

function App(_props: AppRootProps) {
  return (
    <Routes>
      <Route path={ROUTES.Analysis} element={<AnalysisPage />} />
      <Route path="*" element={<AnalysisPage />} />
    </Routes>
  );
}

export default App;

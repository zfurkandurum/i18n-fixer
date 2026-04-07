import React from 'react';
import { useTranslation } from 'react-i18next';

function Header() {
  const { t } = useTranslation();

  return (
    <header>
      <h1>{t('common.save')}</h1>
      <h2>{t('common.cancel')}</h2>
      <h3>{t('common.delete')}</h3>
      <h4>{t('missing.key')}</h4>
      <span>Welcome to our application</span>
      <input placeholder="Search here..." />
    </header>
  );
}

export default Header;

import React from 'react';
import { useTranslation } from 'react-i18next';

function MyComponent() {
  const { t } = useTranslation();

  return (
    <div>
      <h1>{t('dashboard.title')}</h1>
      <p>{t('dashboard.welcome')}</p>
      <button>{t('common.save')}</button>
      <span>Please enter your email address</span>
      <input placeholder="Type here..." />
      <p>Loading data, please wait...</p>
      <a href="https://example.com">Link</a>
      <div className="container">content</div>
    </div>
  );
}

function DynamicExample() {
  const { t } = useTranslation();
  const errorCode = 'network';

  return (
    <div>
      {t(`errors.${errorCode}`)}
      {t('static.key.here')}
    </div>
  );
}

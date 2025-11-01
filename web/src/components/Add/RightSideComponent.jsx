import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { rgba } from 'polished'
import { NoImageIcon } from 'icons'
import {
  FormControl,
  FormHelperText,
  IconButton,
  InputAdornment,
  MenuItem,
  Select,
  TextField,
  useTheme,
  CircularProgress,
  Box,
} from '@material-ui/core'
import { HighlightOff as HighlightOffIcon } from '@material-ui/icons'
import { TORRENT_CATEGORIES } from 'components/categories'

import {
  ClearPosterButton,
  PosterLanguageSwitch,
  RightSide,
  Poster,
  PosterSuggestions,
  PosterSuggestionsItem,
  PosterWrapper,
  RightSideContainer,
} from './style'
import { checkImageURL } from './helpers'
import FileSelector from './FileSelector'

export default function RightSideComponent({
  setTitle,
  setCategory,
  setStrmDir,
  setPosterUrl,
  setIsPosterUrlCorrect,
  setIsUserInteractedWithPoster,
  setPosterList,
  isTorrentSourceCorrect,
  isHashAlreadyExists,
  title,
  category,
  strmDir,
  parsedTitle,
  posterUrl,
  isPosterUrlCorrect,
  posterList,
  currentLang,
  posterSearchLanguage,
  setPosterSearchLanguage,
  posterSearch,
  removePoster,
  torrentSource,
  originalTorrentTitle,
  updateTitleFromSource,
  isCustomTitleEnabled,
  setIsCustomTitleEnabled,
  isEditMode,
  torrentFiles,
  selectedFiles,
  setSelectedFiles,
  showFileSelector,
  setTorrentFiles,
  setShowFileSelector,
  loadingMetadata,
}) {
  const { t } = useTranslation()
  const primary = useTheme().palette.primary.main

  // Для magnet-ссылок НЕ загружаем метаданные автоматически
  // Это будет происходить только после нажатия кнопки Add
  useEffect(() => {
    // Сбрасываем состояние при изменении источника
    if (!torrentSource || !torrentSource.startsWith('magnet:')) {
      setShowFileSelector(false)
      setTorrentFiles([])
      setSelectedFiles([])
    }
  }, [torrentSource, setShowFileSelector, setTorrentFiles, setSelectedFiles])

  const handleTitleChange = ({ target: { value } }) => setTitle(value)
  const handleCategoryChange = ({ target: { value } }) => setCategory(value)
  const handleStrmDirChange = ({ target: { value } }) => setStrmDir(value)
  const handlePosterUrlChange = ({ target: { value } }) => {
    setPosterUrl(value)
    checkImageURL(value).then(setIsPosterUrlCorrect)
    setIsUserInteractedWithPoster(!!value)
    setPosterList()
  }
  const userChangesPosterUrl = url => {
    setPosterUrl(url)
    checkImageURL(url).then(setIsPosterUrlCorrect)
    setIsUserInteractedWithPoster(true)
  }
  // main categories
  const catIndex = TORRENT_CATEGORIES.findIndex(e => e.key === category)
  // const catArray = TORRENT_CATEGORIES.find(e => e.key === category)

  return (
    <RightSide>
      <RightSideContainer isHidden={!isTorrentSourceCorrect || (isHashAlreadyExists && !isEditMode)}>
        {originalTorrentTitle ? (
          <>
            <TextField
              value={originalTorrentTitle}
              margin='dense'
              label={t('AddDialog.OriginalTorrentTitle')}
              style={{ marginTop: '1em' }}
              type='text'
              variant='outlined'
              fullWidth
              disabled={isCustomTitleEnabled}
              InputProps={{ readOnly: true }}
            />
            <TextField
              onChange={handleTitleChange}
              onFocus={() => setIsCustomTitleEnabled(true)}
              onBlur={({ target: { value } }) => !value && setIsCustomTitleEnabled(false)}
              value={title}
              margin='dense'
              label={t('AddDialog.CustomTorrentTitle')}
              type='text'
              variant='outlined'
              fullWidth
              helperText={t('AddDialog.CustomTorrentTitleHelperText')}
              InputProps={{
                endAdornment: (
                  <InputAdornment position='end'>
                    <IconButton
                      size='small'
                      style={{ padding: '1px', marginRight: '-6px' }}
                      onClick={() => {
                        setTitle('')
                        setIsCustomTitleEnabled(!isCustomTitleEnabled)
                        updateTitleFromSource()
                        setIsUserInteractedWithPoster(false)
                      }}
                    >
                      <HighlightOffIcon style={{ color: isCustomTitleEnabled ? primary : rgba('#ccc', 0.25) }} />
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
          </>
        ) : (
          <TextField
            onChange={handleTitleChange}
            value={title}
            margin='dense'
            label={t('AddDialog.TitleBlank')}
            style={{ marginTop: '1em' }}
            type='text'
            variant='outlined'
            fullWidth
            helperText={t('AddDialog.TitleBlankHelperText')}
          />
        )}
        <TextField
          onChange={handlePosterUrlChange}
          value={posterUrl}
          margin='dense'
          label={t('AddDialog.AddPosterLinkInput')}
          type='url'
          variant='outlined'
          fullWidth
        />
        <FormControl fullWidth>
          <FormHelperText style={{ padding: '0.2em 1.2em 0.5em 1.2em' }}>
            {t('AddDialog.CategoryHelperText')}
          </FormHelperText>
          <Select
            labelId='torrent-category-select-label'
            id='torrent-category-select'
            value={category}
            margin='dense'
            onChange={handleCategoryChange}
            variant='outlined'
            fullWidth
            defaultValue=''
            IconComponent={
              category.length > 1
                ? () => (
                    <IconButton
                      size='small'
                      style={{ padding: '1px', marginLeft: '6px', marginRight: '8px' }}
                      onClick={() => {
                        setCategory('')
                      }}
                    >
                      <HighlightOffIcon style={{ color: primary }} />
                    </IconButton>
                  )
                : undefined
            }
          >
            {category.length > 1 && catIndex < 0 ? (
              <MenuItem key={category} value={category}>
                {category}
              </MenuItem>
            ) : (
              ''
            )}

            {TORRENT_CATEGORIES.map(category => (
              <MenuItem key={category.key} value={category.key}>
                {t(category.name)}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        <TextField
          onChange={handleStrmDirChange}
          value={strmDir}
          margin='dense'
          label='Custom .strm Directory'
          type='text'
          variant='outlined'
          fullWidth
          helperText='Optional: Custom folder name for .strm files (between category and torrent folders)'
          placeholder='e.g., Movies2024 or TVShows/Season1'
        />

        {loadingMetadata ? (
          <Box display='flex' flexDirection='column' alignItems='center' padding={3}>
            <CircularProgress />
            <TextField
              margin='dense'
              label={t('LoadingMetadata', 'Loading metadata...')}
              type='text'
              variant='outlined'
              fullWidth
              disabled
              value={t('PleaseWait', 'Please wait while torrent metadata is being downloaded...')}
              helperText={t('ThisMayTakeTime', 'This may take up to 30 seconds')}
              style={{ marginTop: 16 }}
            />
          </Box>
        ) : showFileSelector ? (
          <FileSelector torrentFiles={torrentFiles} selectedFiles={selectedFiles} setSelectedFiles={setSelectedFiles} />
        ) : null}

        <PosterWrapper>
          <Poster poster={+isPosterUrlCorrect}>
            {isPosterUrlCorrect ? <img src={posterUrl} alt='poster' /> : <NoImageIcon />}
          </Poster>

          <PosterSuggestions>
            {posterList
              ?.filter(url => url !== posterUrl)
              .slice(0, 12)
              .map(url => (
                <PosterSuggestionsItem onClick={() => userChangesPosterUrl(url)} key={url}>
                  <img src={url} alt='poster' />
                </PosterSuggestionsItem>
              ))}
          </PosterSuggestions>

          {currentLang !== 'en' && (
            <PosterLanguageSwitch
              onClick={() => {
                const newLanguage = posterSearchLanguage === 'en' ? 'ru' : 'en'
                setPosterSearchLanguage(newLanguage)
                posterSearch(isCustomTitleEnabled ? title : originalTorrentTitle ? parsedTitle : title, newLanguage, {
                  shouldRefreshMainPoster: true,
                })
              }}
              showbutton={+isPosterUrlCorrect}
              color='primary'
              variant='contained'
              size='small'
            >
              {posterSearchLanguage === 'en' ? 'EN' : 'RU'}
            </PosterLanguageSwitch>
          )}

          <ClearPosterButton
            showbutton={+isPosterUrlCorrect}
            onClick={() => {
              removePoster()
              setIsUserInteractedWithPoster(true)
            }}
            color='primary'
            variant='contained'
            size='small'
          >
            {t('Clear')}
          </ClearPosterButton>
        </PosterWrapper>
      </RightSideContainer>

      <RightSideContainer
        isError={torrentSource && (!isTorrentSourceCorrect || (isHashAlreadyExists && !torrentSource.startsWith('magnet:')))}
        notificationMessage={
          !torrentSource
            ? t('AddDialog.AddTorrentSourceNotification')
            : !isTorrentSourceCorrect
            ? t('AddDialog.WrongTorrentSource')
            : isHashAlreadyExists && !torrentSource.startsWith('magnet:') && t('AddDialog.HashExists')
        }
        isHidden={isEditMode || (isTorrentSourceCorrect && (!isHashAlreadyExists || torrentSource.startsWith('magnet:')))}
      />
    </RightSide>
  )
}
